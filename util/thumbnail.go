package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// ExtractThumbnailFromStream 从视频流中截取封面并保存为图片文件
func ExtractThumbnailFromStream(videoStream io.Reader, outputImagePath string, timestamp time.Duration) error {
	// 检查 ffmpeg 是否可用
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	// 创建 ffmpeg 命令
	cmd := exec.Command(
		ffmpegPath,
		"-i", "pipe:0", // 从 stdin 读取输入
		"-ss", fmt.Sprintf("%.2f", timestamp.Seconds()), // 设置时间戳
		"-vframes", "1", // 只提取一帧
		"-f", "image2", // 输出格式为图片
		"-y",            // 覆盖输出文件
		outputImagePath, // 输出图片文件路径
	)

	// 捕获 stderr 输出用于错误诊断
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// 设置 stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	// 将视频流数据复制到 ffmpeg stdin
	go func() {
		defer stdin.Close()
		// 使用带缓冲的复制
		buffer := make([]byte, 32*1024) // 32KB 缓冲区
		_, _ = io.CopyBuffer(stdin, videoStream, buffer)
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v, stderr: %s", err, stderrBuf.String())
	}

	// 检查输出文件是否存在
	if _, err := os.Stat(outputImagePath); os.IsNotExist(err) {
		return fmt.Errorf("failed to create thumbnail, output file does not exist: %s", outputImagePath)
	}

	return nil
}

// ExtractFrameFromStream 从视频流中截取封面并返回图片数据
func ExtractFrameFromStream(videoStream io.Reader, timestamp time.Duration) ([]byte, error) {
	// 检查 ffmpeg 是否可用
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	// 创建 ffmpeg 命令
	cmd := exec.Command(
		ffmpegPath,
		"-i", "pipe:0", // 从 stdin 读取输入
		"-ss", fmt.Sprintf("%.2f", timestamp.Seconds()), // 设置时间戳
		"-vframes", "1", // 只提取一帧
		"-f", "image2", // 输出格式为图片
		"-c:v", "mjpeg", // 使用 JPEG 编码
		"pipe:1", // 输出到 stdout
	)

	// 创建缓冲区来捕获输出
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// 设置 stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	// 将视频流数据复制到 ffmpeg stdin
	go func() {
		defer stdin.Close()
		// 使用带缓冲的复制
		buffer := make([]byte, 32*1024) // 32KB 缓冲区
		_, _ = io.CopyBuffer(stdin, videoStream, buffer)
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg command failed: %v, stderr: %s", err, stderrBuf.String())
	}

	// 返回截取的封面数据
	return stdoutBuf.Bytes(), nil
}

// ExtractThumbnailFromFile 从指定视频文件中截取封面并保存为图片文件（通过流方式处理）
func ExtractThumbnailFromFile(videoFilePath, outputImagePath string, timestamp time.Duration) error {
	// 打开视频文件
	videoFile, err := os.Open(videoFilePath)
	if err != nil {
		return fmt.Errorf("failed to open video file: %v", err)
	}
	defer videoFile.Close()

	// 使用流方式处理
	return ExtractThumbnailFromStream(videoFile, outputImagePath, timestamp)
}
