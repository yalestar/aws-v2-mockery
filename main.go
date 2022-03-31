package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type S3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type S3GetObjectAPI interface {
	GetObject(ctx context.Context, input *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Yapi struct {
	listObjectsAPI S3ListObjectsAPI
	getObjectAPi   S3GetObjectAPI
	client         *s3.Client
}

func (b Yapi) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input,
	optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {

	loo, err := b.client.ListObjectsV2(ctx, input)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return loo, nil
}

func getThemShits(ctx context.Context, api S3ListObjectsAPI,
	bucket string) (*s3.ListObjectsV2Output, error) {
	lsOutput, err := api.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	return lsOutput, nil
}

func main() {
	bucket := "zapp"

	cfg, err := getS3Config("http://localhost:4566")
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := getS3Client(cfg)

	api := Yapi{
		client: client,
	}

	ls, err := getThemShits(context.Background(), api, bucket)
	if err != nil {
		log.Println(err)
	}
	for _, item := range ls.Contents {
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
	}
}
