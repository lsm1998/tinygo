package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minIOClient struct {
	*minio.Client
}

func NewMinIOClient(endpoint, accessKeyID, secretAccessKey string, secure bool) (*minIOClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	return &minIOClient{Client: minioClient}, nil
}
