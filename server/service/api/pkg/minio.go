package pkg

import (
	"GoYin/server/service/api/config"
	"context"
	"github.com/minio/minio-go/v7"
	"os"
)

func minioUpgrade() {
	client := config.GlobalMinioClient
	ctx := context.Background()

	videoFilePath := "path/to/video.mp4"
	videoObjectName := "videos/video.mp4"

	coverImagePath := "path/to/cover.jpg"
	coverImageObjectName := "covers/cover.jpg"

	// 上传视频
	videoFile, err := os.Open(videoFilePath)
	if err != nil {
		// 处理文件打开错误
	}
	defer videoFile.Close()

	videoStat, err := videoFile.Stat()
	if err != nil {
		// 处理获取文件信息错误
	}

	_, err = client.PutObject(ctx, "bucket_name", videoObjectName, videoFile, videoStat.Size(), minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		// 处理上传错误
	}

	// 上传封面图片
	coverImageFile, err := os.Open(coverImagePath)
	if err != nil {
		// 处理文件打开错误
	}
	defer coverImageFile.Close()

	coverImageStat, err := coverImageFile.Stat()
	if err != nil {
		// 处理获取文件信息错误
	}

	_, err = client.PutObject(ctx, "bucket_name", coverImageObjectName, coverImageFile, coverImageStat.Size(), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		// 处理上传错误
	}
}
