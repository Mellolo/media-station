package bizVideo

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Mellolo/common/errors"
	"github.com/Mellolo/media-station/models/dto/contextDTO"
	"github.com/Mellolo/media-station/utils/videoUtil"
	"github.com/beego/beego/v2/core/logs"
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
}

// 全局 HLS 会话管理器
var hlsSessionManager = &HLSSessionManager{
	sessions: make(map[string]*HLSSession),
}

// StreamVideoToHLS 流式转码入口
func (impl *VideoBizServiceImpl) StreamVideoToHLS(ctx contextDTO.ContextDTO, id int64) map[string]string {
	video, err := impl.videoMapper.SelectById(id)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("video [%d] doesn't exist", id)))
	}

	// 判断是否需要转码（从 video_url 提取扩展名）
	ext := videoUtil.GetFileExtension(video.VideoUrl)
	if !videoUtil.NeedsTranscoding(ext) {
		// 浏览器原生支持，直接返回原始文件地址
		return map[string]string{
			"type": "native",
			"url":  fmt.Sprintf("/api/video/play/%d", id),
		}
	}

	// 需要转码，创建唯一会话ID
	sessionID := fmt.Sprintf("%d_%d", id, time.Now().UnixNano())
	hlsDir := fmt.Sprintf("/tmp/hls_streaming/%s", sessionID)

	// 创建临时目录
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		panic(errors.WrapError(err, "create hls directory failed"))
	}

	// 获取原始文件的预签名 URL
	inputURL := impl.videoStorage.GetStreamURL(bucketVideo, video.VideoUrl, 30*time.Minute)

	playlistPath := fmt.Sprintf("%s/playlist.m3u8", hlsDir)
	segmentPattern := fmt.Sprintf("%s/segment_%03d.ts", hlsDir)

	// 启动 FFmpeg 后台转码
	cmd := exec.Command("ffmpeg",
		"-i", inputURL,
		"-c:v", "libx264",
		"-preset", "ultrafast", // 最快编码速度，降低延迟
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", "2", // 2秒片段（更快起播）
		"-hls_list_size", "0", // 保留所有片段
		"-hls_segment_filename", segmentPattern,
		"-pix_fmt", "yuv420p",
		playlistPath)

	// 启动进程（异步，不等待完成）
	if err := cmd.Start(); err != nil {
		os.RemoveAll(hlsDir)
		panic(errors.WrapError(err, "start ffmpeg failed"))
	}

	// 注册会话（用于后续清理）
	hlsSessionManager.Register(sessionID, hlsDir, cmd.Process)

	logs.Info(fmt.Sprintf("HLS转码已启动: session=%s, video=%d", sessionID, id))

	// 立即返回 playlist URL（前端可以开始请求）
	return map[string]string{
		"type":         "hls",
		"playlist_url": fmt.Sprintf("/api/video/hls/%s/playlist.m3u8", sessionID),
		"session_id":   sessionID,
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
func (m *HLSSessionManager) Register(sessionID, dir string, process *os.Process) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sessions[sessionID] = &HLSSession{
		SessionID: sessionID,
		Directory: dir,
		Process:   process,
		CreatedAt: time.Now(),
	}

	// 启动清理协程（10分钟后清理）
	go func() {
		time.Sleep(10 * time.Minute)
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

// GetFileExtension 从文件路径提取扩展名
func GetFileExtension(filePath string) string {
	parts := strings.Split(filePath, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// NeedsTranscoding 判断视频格式是否需要转码
func NeedsTranscoding(fileExtension string) bool {
	// 浏览器原生支持的格式
	supportedExts := map[string]bool{
		".mp4":  true, // H.264/AAC
		".webm": true, // VP8/VP9/Vorbis
		".ogv":  true, // Theora/Vorbis
	}

	ext := strings.ToLower(fileExtension)
	// 确保扩展名以 . 开头
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return !supportedExts[ext]
}

// GetExtensionFromFormat 从 ffprobe 格式名映射到文件扩展名
func GetExtensionFromFormat(formatName string) string {
	// 处理可能包含多个格式名的情况（如 "mov,mp4,m4a"）
	formatParts := strings.Split(formatName, ",")
	if len(formatParts) > 0 {
		formatName = formatParts[0]
	}

	formatMap := map[string]string{
		"matroska":  "mkv",
		"avi":       "avi",
		"mov":       "mp4",
		"mp4":       "mp4",
		"flv":       "flv",
		"rm":        "rmvb",
		"realmedia": "rmvb",
		"mpeg":      "mpg",
		"mpegts":    "ts",
		"webm":      "webm",
		"asf":       "wmv",
		"wmv":       "wmv",
		"quicktime": "mov",
		"ogg":       "ogv",
	}

	if ext, ok := formatMap[formatName]; ok {
		return ext
	}

	// 默认返回 mp4
	return "mp4"
}
