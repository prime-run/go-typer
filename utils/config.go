package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const configDirName = "go-typer"

func GetConfigDirPath() string {
	configDir, err := GetAppConfigDir()
	if err != nil {
		return os.TempDir()
	}
	return configDir
}

func GetAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, configDirName)

	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return appConfigDir, nil
}
