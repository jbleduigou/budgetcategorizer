package config

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbleduigou/budgetcategorizer/iface"
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
func GetConfiguration(ctx context.Context, cli iface.S3DownloadAPI) Configuration {
	bucket, ok := os.LookupEnv("CONFIGURATION_FILE_BUCKET")
	if !ok {
		return defaultConfig()
	}
	objectKey, ok := os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
	if !ok {
		return defaultConfig()
	}
	slog.Info("Downloading configuration file from S3",
		slog.String("bucket", bucket),
		slog.String("object-key", objectKey))
	output, err := cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		slog.Error("Error while downloading configuration file",
			slog.String("bucket", bucket),
			slog.String("object-key", objectKey),
			slog.Any("error", err))
		return defaultConfig()
	}
	defer output.Body.Close()
	body, err := io.ReadAll(output.Body)
	if err != nil {
		slog.Error("Error while downloading configuration file",
			slog.String("bucket", bucket),
			slog.String("object-key", objectKey),
			slog.Any("error", err))
		return defaultConfig()
	}
	slog.Info("Successfully downloaded configuration file from S3",
		slog.String("bucket", bucket),
		slog.String("object-key", objectKey))
	return parseConfiguration(body)
}

func defaultConfig() Configuration {
	slog.Warn("Using default configuration")
	return parseConfiguration([]byte(defaultConfiguration))
}

func parseConfiguration(yml []byte) Configuration {
	c := Configuration{}
	_ = yaml.Unmarshal(yml, &c)
	for _, v := range c.Categories {
		c.Keywords[v] = v
	}
	return c
}
