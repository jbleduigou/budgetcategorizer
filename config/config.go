package config

import (
	"fmt"
	"os"
	"strings"

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

// Configuration provides and interface for the configuration of the software
type Configuration struct {
	Categories []string
	Keywords   map[string]string
}

// GetConfiguration will return the configuration.
// Configuration can be downloaded from an S3 bucket.
// If not available a default configuration will be returned
func GetConfiguration(downloader s3manageriface.DownloaderAPI, requestID string) Configuration {
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
				fmt.Printf("[INFO] %v Downloaded configuration file '%s' from bucket '%s' \n", requestID, objectKey, bucket)
				return parseConfiguration([]byte(buff.Bytes()))
			}
			fmt.Printf("[WARN] %v Could not download configuration file '%s' from bucket '%s'\n", requestID, objectKey, bucket)
			fmt.Printf("[WARN] %v %v\n", requestID, strings.ReplaceAll(err.Error(), "\n", " "))
		}
	}
	fmt.Printf("[WARN] %v Using default configuration \n", requestID)
	return parseConfiguration([]byte(defaultConfiguration))
}

func parseConfiguration(yml []byte) Configuration {
	c := Configuration{}
	yaml.Unmarshal(yml, &c)
	for _, v := range c.Categories {
		c.Keywords[v] = v
	}
	return c
}
