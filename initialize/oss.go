package initialize

import (
	"github.com/Mellolo/common/oss"
	"github.com/beego/beego/v2/server/web"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitOss() {
	client, err := minio.NewCore(
		web.AppConfig.DefaultString("minio::minio_end_point", "localhost:9000"),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				web.AppConfig.DefaultString("minio::minio_access_key", "minioadmin"),
				web.AppConfig.DefaultString("minio::minio_secret_key", "minioadmin"),
				""),
			Secure: false,
		})
	if err != nil {
		panic(err)
	}
	err = oss.InitOss(oss.NewMinioClient(client))
	if err != nil {
		panic(err)
	}
}
