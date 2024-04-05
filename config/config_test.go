package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

type mockReadCloser struct{}

func (m *mockReadCloser) Close() error {
	return nil
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestGetConfigurationShouldUseDefaultGivenNoBucket(t *testing.T) {
	os.Unsetenv("CONFIGURATION_FILE_BUCKET")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")

	m := mock.NewDownloader("")
	configuration := GetConfiguration(context.Background(), m)

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "Courses Alimentation")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Express Proxi Saint Thonan"], "Courses Alimentation")
	assert.Equal(t, configuration.Keywords["Courses Alimentation"], "Courses Alimentation")
	m.AssertExpectations(t)
}

func TestGetConfigurationShouldUseDefaultGivenNoObjectKey(t *testing.T) {
	os.Setenv("CONFIGURATION_FILE_BUCKET", "mybucket")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")

	m := mock.NewDownloader("")
	configuration := GetConfiguration(context.Background(), m)

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "Courses Alimentation")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Express Proxi Saint Thonan"], "Courses Alimentation")
	assert.Equal(t, configuration.Keywords["Courses Alimentation"], "Courses Alimentation")
	m.AssertExpectations(t)
}

func TestGetConfigurationShouldUseDefaultGivenErrorWithS3(t *testing.T) {
	yamlContent := "categories:\n  - MyCategory\nkeywords:\n  Tesco London: MyCategory"

	os.Setenv("CONFIGURATION_FILE_BUCKET", "mybucket")
	os.Setenv("CONFIGURATION_FILE_OBJECT_KEY", "configuration.yaml")
	m := mock.NewDownloader(yamlContent)
	m.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("configuration.yaml")},
		mock.Anything).Return(int64(1337), fmt.Errorf("error for unit test"))
	configuration := GetConfiguration(context.Background(), m)

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "Courses Alimentation")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Express Proxi Saint Thonan"], "Courses Alimentation")
	assert.Equal(t, configuration.Keywords["Courses Alimentation"], "Courses Alimentation")

	os.Unsetenv("CONFIGURATION_FILE_BUCKET")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")
}

func TestGetConfigurationShouldUseDefaultGivenErrorWithReader(t *testing.T) {
	yamlContent := "categories:\n  - MyCategory\nkeywords:\n  Tesco London: MyCategory"

	os.Setenv("CONFIGURATION_FILE_BUCKET", "mybucket")
	os.Setenv("CONFIGURATION_FILE_OBJECT_KEY", "configuration.yaml")
	m := mock.NewDownloader(yamlContent)
	m.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("configuration.yaml")},
		mock.Anything).Return(&s3.GetObjectOutput{Body: &mockReadCloser{}}, nil)
	configuration := GetConfiguration(context.Background(), m)

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "Courses Alimentation")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Express Proxi Saint Thonan"], "Courses Alimentation")
	assert.Equal(t, configuration.Keywords["Courses Alimentation"], "Courses Alimentation")
	m.AssertExpectations(t)

	os.Unsetenv("CONFIGURATION_FILE_BUCKET")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")
}

func TestGetConfigurationShouldDownload(t *testing.T) {
	yamlContent := "categories:\n  - MyCategory\nkeywords:\n  Tesco London: MyCategory"

	os.Setenv("CONFIGURATION_FILE_BUCKET", "mybucket")
	os.Setenv("CONFIGURATION_FILE_OBJECT_KEY", "configuration.yaml")
	m := mock.NewDownloader(yamlContent)
	m.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("configuration.yaml")},
		mock.Anything).Return(&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(yamlContent))}, nil)
	configuration := GetConfiguration(context.Background(), m)

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "MyCategory")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Tesco London"], "MyCategory")
	assert.Equal(t, configuration.Keywords["MyCategory"], "MyCategory")
	m.AssertExpectations(t)

	os.Unsetenv("CONFIGURATION_FILE_BUCKET")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")
}
