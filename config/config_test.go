package config

import (
	"fmt"
	"testing"

	"github.com/jbleduigou/budgetcategorizer/mock"
)

func TestGetConfigurationShouldUseDefault(t *testing.T) {
	configuration := GetConfiguration(mock.NewDownloader(""))

	fmt.Println(configuration)
}

func TestGetConfigurationShouldDownload(t *testing.T) {
	configuration := GetConfiguration(mock.NewDownloader(""))

	fmt.Println(configuration)
}
