package initialize

import (
	"GoYin/server/service/api/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() {
	s3Client, err := minio.New(config.GlobalServerConfig.MinioInfo.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GlobalServerConfig.MinioInfo.AccessKeyID, config.GlobalServerConfig.MinioInfo.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		hlog.Fatal(err)
	}
	config.GlobalMinioClient = s3Client
}
