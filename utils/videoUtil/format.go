package videoUtil

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// VideoFormatInfo 视频格式信息结构体
type VideoFormatInfo struct {
	FormatName string  `json:"format_name"` // matroska, avi, mov, rm
	MimeType   string  `json:"mime_type"`   // video/x-matroska, video/avi
	VideoCodec string  `json:"video_codec"` // h264, rv40 (for rmvb)
	AudioCodec string  `json:"audio_codec"` // aac, mp3
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Bitrate    int     `json:"bitrate"`
	Duration   float64 `json:"duration"`
}

// DetectVideoFormat 使用 ffprobe 检测视频格式
func DetectVideoFormat(filePath string) (*VideoFormatInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe执行失败: %v", err)
	}

	// 解析 JSON 输出
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	info := &VideoFormatInfo{}

	// 提取 format 信息
	if format, ok := result["format"].(map[string]interface{}); ok {
		if formatName, ok := format["format_name"].(string); ok {
			info.FormatName = formatName
		}
		if mimeType, ok := format["mime_type"].(string); ok {
			info.MimeType = mimeType
		}
		if duration, ok := format["duration"].(string); ok {
			info.Duration, _ = strconv.ParseFloat(duration, 64)
		}
		if bitrate, ok := format["bit_rate"].(string); ok {
			info.Bitrate, _ = strconv.Atoi(bitrate)
		}
	}

	// 提取 streams 信息
	if streams, ok := result["streams"].([]interface{}); ok {
		for _, stream := range streams {
			s, ok := stream.(map[string]interface{})
			if !ok {
				continue
			}

			codecType, ok := s["codec_type"].(string)
			if !ok {
				continue
			}

			if codecType == "video" && info.VideoCodec == "" {
				if codecName, ok := s["codec_name"].(string); ok {
					info.VideoCodec = codecName
				}
				if width, ok := s["width"].(float64); ok {
					info.Width = int(width)
				}
				if height, ok := s["height"].(float64); ok {
					info.Height = int(height)
				}
			} else if codecType == "audio" && info.AudioCodec == "" {
				if codecName, ok := s["codec_name"].(string); ok {
					info.AudioCodec = codecName
				}
			}
		}
	}

	return info, nil
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

// GetFileExtension 从文件路径提取扩展名
func GetFileExtension(filePath string) string {
	parts := strings.Split(filePath, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// GetContentTypeByExtension 根据文件扩展名获取 MIME Content-Type
func GetContentTypeByExtension(filePath string) string {
	lower := strings.ToLower(filePath)
	switch {
	case strings.HasSuffix(lower, ".mp4") || strings.HasSuffix(lower, ".m4v"):
		return "video/mp4"
	case strings.HasSuffix(lower, ".webm"):
		return "video/webm"
	case strings.HasSuffix(lower, ".mkv"):
		return "video/x-matroska"
	case strings.HasSuffix(lower, ".avi"):
		return "video/x-msvideo"
	case strings.HasSuffix(lower, ".mov"):
		return "video/quicktime"
	case strings.HasSuffix(lower, ".flv"):
		return "video/x-flv"
	case strings.HasSuffix(lower, ".wmv"):
		return "video/x-ms-wmv"
	case strings.HasSuffix(lower, ".ts"):
		return "video/mp2t"
	case strings.HasSuffix(lower, ".ogv") || strings.HasSuffix(lower, ".ogg"):
		return "video/ogg"
	case strings.HasSuffix(lower, ".mpg") || strings.HasSuffix(lower, ".mpeg"):
		return "video/mpeg"
	case strings.HasSuffix(lower, ".3gp"):
		return "video/3gpp"
	case strings.HasSuffix(lower, ".rmvb") || strings.HasSuffix(lower, ".rm"):
		return "application/vnd.rn-realmedia"
	case strings.HasSuffix(lower, ".asf"):
		return "video/x-ms-asf"
	default:
		return "video/mp4"
	}
}
