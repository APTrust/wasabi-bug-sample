package main

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go name-of-bucket")
		return
	}
	bucket := os.Args[1]
	client := getClient()

	// Run this test twice, so we can see the first upload
	// always succeeds and the second always fails.
	doUploads(client, bucket)
	doUploads(client, bucket)
}

func doUploads(client *minio.Client, bucket string) {
	goodPutOptions := getPutOptions("Metadata values with single spaces are OK")
	badPutOptions := getPutOptions("Metadata values with two consecutive  spaces cause upload to fail")
	err := uploadFile(client, bucket, goodPutOptions)
	if err != nil {
		fmt.Println("Upload with goodPutOptions FAILED with error:", err)
	} else {
		fmt.Println("Upload with goodPutOptions SUCCEEDED")
	}

	err = uploadFile(client, bucket, badPutOptions)
	if err != nil {
		fmt.Println("Upload with badPutOptions FAILED with error:", err)
	} else {
		fmt.Println("Upload with badPutOptions SUCCEEDED")
	}
}

func getPutOptions(str string) minio.PutObjectOptions {
	return minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"custom-data": str,
		},
		ContentType: "text/plain",
	}
}

func getClient() *minio.Client {
	client, err := minio.New(
		"s3.us-east-1.wasabisys.com",
		&minio.Options{
			Creds:  credentials.NewStaticV4(getEnvVar("WASABI_ACCESS_KEY"), getEnvVar("WASABI_SECRET_KEY"), ""),
			Secure: true,
		})
	if err != nil {
		panic(err)
	}
	return client
}

func getEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("Env var %s is not set", name))
	}
	return value
}

func uploadFile(client *minio.Client, bucket string, putOptions minio.PutObjectOptions) error {
	file, err := os.Open("sample.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = client.PutObject(
		context.Background(),
		bucket,
		"sample.txt",
		file,
		342,
		putOptions,
	)
	return err
}
