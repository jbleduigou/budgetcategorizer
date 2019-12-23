package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"gopkg.in/yaml.v2"
)

const defaultConfiguration = `
categories:
  - Courses Alimentation
keywords:
  Express Proxi Saint Thonan: Courses Alimentation
`

type Configuration struct {
	Categories []string
	Keywords   map[string]string
}

func GetConfiguration(downloader s3manageriface.DownloaderAPI) Configuration {
	bucket, ok := os.LookupEnv("CONFIGURATION_FILE_BUCKET")
	if ok {
		objectKey, ok := os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
		if ok {
			buff := &aws.WriteAtBuffer{}
			_, err := downloader.Download(buff, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(objectKey),
			})
			if err == nil {
				fmt.Printf("Download configuration file '%s' from bucket '%s' \n", objectKey, bucket)
				return parseConfiguration([]byte(buff.Bytes()))
			}
		}
	}
	return parseConfiguration([]byte(defaultConfiguration))
}

func parseConfiguration(yml []byte) Configuration {
	c := Configuration{}
	yaml.Unmarshal(yml, &c)
	return c
}
