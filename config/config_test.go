package config

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigurationShouldUseDefault(t *testing.T) {
	m := mock.NewDownloader("")
	configuration := GetConfiguration(m, "1a5931ca-dd5d-11ea-90cb-3822e2348205")

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "Courses Alimentation")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Express Proxi Saint Thonan"], "Courses Alimentation")
	assert.Equal(t, configuration.Keywords["Courses Alimentation"], "Courses Alimentation")
	m.AssertExpectations(t)
}

func TestGetConfigurationShouldDownload(t *testing.T) {
	yamlContent := "categories:\n  - MyCategory\nkeywords:\n  Tesco London: MyCategory"

	os.Setenv("CONFIGURATION_FILE_BUCKET", "mybucket")
	os.Setenv("CONFIGURATION_FILE_OBJECT_KEY", "configuration.yaml")
	m := mock.NewDownloader(yamlContent)
	m.On("Download",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("configuration.yaml")},
		mock.Anything).Return(int64(1337), nil)
	configuration := GetConfiguration(m, "1a5931ca-dd5d-11ea-90cb-3822e2348205")

	assert.Equal(t, len(configuration.Categories), 1)
	assert.Equal(t, configuration.Categories[0], "MyCategory")
	assert.Equal(t, len(configuration.Keywords), 2)
	assert.Equal(t, configuration.Keywords["Tesco London"], "MyCategory")
	assert.Equal(t, configuration.Keywords["MyCategory"], "MyCategory")
	m.AssertExpectations(t)

	os.Unsetenv("CONFIGURATION_FILE_BUCKET")
	os.Unsetenv("CONFIGURATION_FILE_OBJECT_KEY")
}
