package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var readDirFlag, s3bucketARN, s3bucketName, s3bucketPath, accessKeyId, secretAccessKey, awsRegion string

func getFolderContentList(folder2read string) []string {
	var files []string

	err := filepath.Walk(folder2read, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func createConnector() *s3.S3 {
	return s3.New(session.New())
}

func uploadFiles(connector *s3.S3, files []string) {
	for _, file := range files {
		input := &s3.UploadPartInput{
			Body:                aws.ReadSeekCloser(strings.NewReader(file)),
			Bucket:              aws.String(s3bucketName),
			Key:                 aws.String(accessKeyId),
			ExpectedBucketOwner: aws.String(accessKeyId),
			// PartNumber: aws.Int64(1),
		}

		result, err := connector.UploadPart(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println(result)
	}
}

func main() {
	readFiles := getFolderContentList(readDirFlag)

	fmt.Println(len(readFiles))
	fmt.Println(readDirFlag)
	fmt.Println(s3bucketARN)
	fmt.Println(s3bucketPath)
}

func init() {
	gotEnvError := godotenv.Load()

	if gotEnvError != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&readDirFlag, "readDir", os.Getenv("L2S3_READ_DIR"), "Set a local directory path to analyse")
	flag.StringVar(&s3bucketARN, "bucketARN", os.Getenv("L2S3_USE_BUCKET"), "Set a remote S3 bucket ARN")
	flag.StringVar(&s3bucketPath, "bucketPath", os.Getenv("L2S3_BUCKET_PATH"), "Set the remote S3 bucket folder path")

	flag.Parse()
}
