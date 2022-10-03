package common

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func testDumpToFile(path string, content string) error {
	jsonFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	if _, err := jsonFile.WriteString(content); err != nil {
		return err
	}
	return nil
}

func testAccessConfig() *AccessConfig {
	return &AccessConfig{}
}

func TestGetconfigFromFile(t *testing.T) {
	if err := os.Setenv("OSC_PROFILE", "SomeProfile"); err != nil {
		t.Fatalf("Cannot set OSC_PROFILE: %s", err.Error())
	}
	ak := os.Getenv("OSC_ACCESS_KEY")
	sk := os.Getenv("OSC_SECRET_KEY")
	region := os.Getenv("OSC_REGION")
	os.Unsetenv("OSC_ACCESS_KEY")
	os.Unsetenv("OSC_SECRET_KEY")
	os.Unsetenv("OSC_REGION")
	defer os.Setenv("OSC_ACCESS_KEY", ak)
	defer os.Setenv("OSC_SECRET_KEY", sk)
	defer os.Setenv("OSC_REGION", region)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Cannot get user home dir: %s", err.Error())
	}
	configFolderPath := path.Join(home, ".osc")
	configPath := path.Join(configFolderPath, "config.json")

	os.RemoveAll(configFolderPath)
	if err := os.Mkdir(configFolderPath, os.ModePerm); err != nil {
		t.Fatalf("Cannot create .osc folder: %s", err.Error())
	}
	defer os.RemoveAll(configFolderPath)

	content := fmt.Sprintf(`{
		"SomeProfile": {
			"access_key": "%s",
			"secret_key": "%s",
			"endpoints": {
				"api": "api.%s.outscale.com/api/v1"
			}
		}}`, ak, sk, region)
	if err := testDumpToFile(configPath, content); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	defer os.Remove(configPath)
	accessConfig := &AccessConfig{}
	_, err = accessConfig.NewOSCClient()
	if err != nil {
		t.Fatalf("Cannot create accessConfig: %s", err.Error())
	}
}
