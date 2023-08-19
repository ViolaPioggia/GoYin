package pkg

import (
	"GoYin/server/service/api/config"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/minio/minio-go/v7"
)

func MinioUpgrade(suffix string, tmpFilePath string, fileName string) error {
	res, err := config.GlobalMinioClient.FPutObject(context.Background(), config.GlobalServerConfig.MinioInfo.Bucket, fileName, tmpFilePath, minio.PutObjectOptions{
		ContentType: "application/" + suffix,
	})
	fmt.Println(res)
	if err != nil {
		hlog.Error("minio upgrade failed,", err)
		return err
	}
	return nil
}
