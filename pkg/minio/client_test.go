package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123456"
	client, err := NewMinIOClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v\n", client) // minioClient is now set up
	exists, err := client.BucketExists(context.Background(), "yidu")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(exists)
}

func TestPutObject(t *testing.T) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123456"
	client, err := NewMinIOClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("D:\\图片\\2016\\0a4e11921c69fee3002a0c7dafca9991.jpg")
	if err != nil {
		t.Fatal(err)
	}
	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	info, err := client.PutObject(context.Background(), "yidu", "demo02.jpg", file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}

func TestGetObject(t *testing.T) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123456"
	client, err := NewMinIOClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		t.Fatal(err)
	}
	object, err := client.GetObject(context.Background(), "yidu", "demo01.jpg", minio.GetObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	info, err := object.Stat()
	fmt.Println(info)
}

func TestNotification(t *testing.T) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123456"
	client, err := NewMinIOClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		t.Fatal(err)
	}
	notification := client.ListenBucketNotification(context.Background(), "yidu", "", "", []string{
		"s3:ObjectCreated:*", "s3:ObjectRemoved:*",
	})
	for {
		ev := <-notification
		if ev.Err != nil {
			t.Error(ev.Err)
			return
		}
		for _, v := range ev.Records {
			fmt.Println(v.EventName)
		}
	}
}
