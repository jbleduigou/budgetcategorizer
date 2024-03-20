package config

import (
	"log/slog"
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

// Configuration provides and interface for the configuration of the software
type Configuration struct {
	Categories []string
	Keywords   map[string]string
}

// GetConfiguration will return the configuration.
// Configuration can be downloaded from an S3 bucket.
// If not available a default configuration will be returned
func GetConfiguration(downloader s3manageriface.DownloaderAPI) Configuration {
	bucket, ok := os.LookupEnv("CONFIGURATION_FILE_BUCKET")
	if ok {
		objectKey, ok := os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
		if ok {
			slog.Info("Downloading configuration file from S3",
				slog.String("bucket", bucket),
				slog.String("object-key", objectKey))
			buff := &aws.WriteAtBuffer{}
			_, err := downloader.Download(buff, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(objectKey),
			})
			if err == nil {
				slog.Info("Successfully downloaded configuration file from S3",
					slog.String("bucket", bucket),
					slog.String("object-key", objectKey))
				return parseConfiguration([]byte(buff.Bytes()))
			}
			slog.Error("Error while downloading configuration file",
				slog.String("bucket", bucket),
				slog.String("object-key", objectKey),
				slog.Any("error", err))
		}
	}
	slog.Warn("Using default configuration")
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
