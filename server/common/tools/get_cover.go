package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetVideoCover(videoPath string) (string, error) {
	// 获取视频文件名（不包含扩展名）
	videoFileName := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))

	// 生成封面图片文件名
	coverFileName := videoFileName + ".jpg"

	// 构造 FFmpeg 命令
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01", "-vframes", "1", coverFileName)
	cmd.Stderr = os.Stderr

	// 执行 FFmpeg 命令
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return coverFileName, nil
}
