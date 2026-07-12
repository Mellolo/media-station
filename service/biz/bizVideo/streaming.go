package bizVideo

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/common/utils/videoUtil"
	"github.com/Mellolo/media-station/models/dto/contextDTO"
	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
)

// HLSSessionManager HLS 会话管理器
type HLSSessionManager struct {
	sessions map[string]*HLSSession
	mutex    sync.RWMutex
}

// HLSSession HLS 会话信息
type HLSSession struct {
	SessionID string
	Directory string
	Process   *os.Process
	CreatedAt time.Time
	VideoID   int64
}

// 全局 HLS 会话管理器
var hlsSessionManager = &HLSSessionManager{
	sessions: make(map[string]*HLSSession),
}

// StreamVideoToHLS 流式转码入口
// 判断视频格式是否需要转码：
// - 浏览器原生格式(mp4/webm/ogv) → 直接返回原生播放地址
// - 非原生格式(mkv/avi/rmvb等) → 启动 FFmpeg 实时转码为 HLS
func (impl *VideoBizServiceImpl) StreamVideoToHLS(ctx contextDTO.ContextDTO, id int64) map[string]string {
	video, err := impl.videoMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}

	ext := videoUtil.GetFileExtension(video.VideoUrl)

	// 浏览器原生支持的格式，直接播放
	if !videoUtil.NeedsTranscoding(ext) {
		contentType := videoUtil.GetContentTypeByExtension(video.VideoUrl)
		return map[string]string{
			"type":         "native",
			"url":          fmt.Sprintf("/api/video/play/%d", id),
			"content_type": contentType,
		}
	}

	// 需要转码：启动 FFmpeg 实时转码为 HLS
	logs.Info(fmt.Sprintf("视频 [%d] 格式 [%s] 需要转码为 HLS", id, ext))

	// 检查是否已有活跃的 HLS 会话（避免重复转码）
	hlsSessionManager.mutex.RLock()
	for _, session := range hlsSessionManager.sessions {
		if session.VideoID == id {
			hlsSessionManager.mutex.RUnlock()
			logs.Info(fmt.Sprintf("视频 [%d] 已有活跃 HLS 会话: %s", id, session.SessionID))
			return map[string]string{
				"type":         "hls",
				"playlist_url": fmt.Sprintf("/api/video/hls/%s/playlist.m3u8", session.SessionID),
				"content_type": "application/x-mpegURL",
			}
		}
	}
	hlsSessionManager.mutex.RUnlock()

	// 生成会话 ID
	sessionID := uuid.New().String()

	// 创建 HLS 输出目录
	hlsDir := filepath.Join("/tmp/hls_streaming", sessionID)
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		panic(errors.WrapError(err, "创建 HLS 输出目录失败"))
	}

	// 获取 MinIO 预签名 URL 作为 FFmpeg 输入
	inputURL := impl.videoStorage.GetStreamURL(bucketVideo, video.VideoUrl, 30*time.Minute)

	// 构建 FFmpeg 命令
	outputPattern := filepath.Join(hlsDir, "segment_%03d.ts")
	playlistPath := filepath.Join(hlsDir, "playlist.m3u8")

	cmd := exec.Command("ffmpeg",
		"-i", inputURL,
		"-c:v", "libx264", // 视频编码为 H.264
		"-c:a", "aac", // 音频编码为 AAC
		"-preset", "veryfast", // 编码速度优先
		"-tune", "zerolatency", // 低延迟优化
		"-start_number", "0",
		"-hls_time", "6", // 每个片段 6 秒
		"-hls_list_size", "0", // 保留所有片段
		"-hls_segment_filename", outputPattern,
		"-f", "hls",
		playlistPath,
	)

	// 启动 FFmpeg 进程
	if err := cmd.Start(); err != nil {
		os.RemoveAll(hlsDir)
		panic(errors.WrapError(err, "启动 FFmpeg 转码失败"))
	}

	logs.Info(fmt.Sprintf("FFmpeg 转码已启动: video=%d, session=%s, pid=%d", id, sessionID, cmd.Process.Pid))

	// 注册会话（10 分钟后自动清理）
	hlsSessionManager.Register(sessionID, hlsDir, cmd.Process, id)

	// 启动后台等待 FFmpeg 进程结束
	go func() {
		if err := cmd.Wait(); err != nil {
			logs.Warn(fmt.Sprintf("FFmpeg 转码进程结束: session=%s, err=%v", sessionID, err))
		} else {
			logs.Info(fmt.Sprintf("FFmpeg 转码完成: session=%s", sessionID))
		}
	}()

	return map[string]string{
		"type":         "hls",
		"playlist_url": fmt.Sprintf("/api/video/hls/%s/playlist.m3u8", sessionID),
		"content_type": "application/x-mpegURL",
	}
}

// ServeHLSSegment 提供 HLS 片段服务
func (impl *VideoBizServiceImpl) ServeHLSSegment(sessionID string, filename string) (string, error) {
	filePath := fmt.Sprintf("/tmp/hls_streaming/%s/%s", sessionID, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果文件还未生成，等待一段时间（最多5秒）
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			if _, err := os.Stat(filePath); err == nil {
				break
			}
		}

		// 仍然不存在，返回错误
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", errors.NewError("HLS segment not ready")
		}
	}

	return filePath, nil
}

// Register 注册 HLS 会话
func (m *HLSSessionManager) Register(sessionID, dir string, process *os.Process, videoID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sessions[sessionID] = &HLSSession{
		SessionID: sessionID,
		Directory: dir,
		Process:   process,
		CreatedAt: time.Now(),
		VideoID:   videoID,
	}

	// 启动清理协程（30 分钟后清理，给转码足够时间）
	go func() {
		time.Sleep(30 * time.Minute)
		m.Cleanup(sessionID)
	}()
}

// Cleanup 清理 HLS 会话
func (m *HLSSessionManager) Cleanup(sessionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		logs.Info(fmt.Sprintf("清理HLS会话: session=%s", sessionID))

		// 终止 FFmpeg 进程
		if session.Process != nil {
			session.Process.Kill()
		}

		// 删除临时目录
		if err := os.RemoveAll(session.Directory); err != nil {
			logs.Error(fmt.Sprintf("删除临时目录失败: %v", err))
		}

		// 移除会话记录
		delete(m.sessions, sessionID)
	}
}

// GetSession 获取会话信息
func (m *HLSSessionManager) GetSession(sessionID string) (*HLSSession, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[sessionID]
	return session, exists
}

// ServeHLSFile 提供 HLS 文件服务（HTTP handler）
func ServeHLSFile(filePath string, w http.ResponseWriter, r *http.Request) {
	// 设置正确的 Content-Type
	if strings.Contains(filePath, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.Contains(filePath, ".ts") {
		w.Header().Set("Content-Type", "video/MP2T")
	}

	// 提供静态文件服务
	http.ServeFile(w, r, filePath)
}
