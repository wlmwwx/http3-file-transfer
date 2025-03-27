package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Server struct {
		Host         string   `yaml:"host"`
		Port         int      `yaml:"port"`
		CertFile     string   `yaml:"cert_file"`
		KeyFile      string   `yaml:"key_file"`
		UploadDir    string   `yaml:"upload_dir"`
		MaxFileSize  int64    `yaml:"max_file_size"`
		AllowedTypes []string `yaml:"allowed_types"`
	} `yaml:"server"`
}

type ClientConfig struct {
	Client struct {
		ServerHost  string `yaml:"server_host"`
		ServerPort  int    `yaml:"server_port"`
		DownloadDir string `yaml:"download_dir"`
		MaxRetries  int    `yaml:"max_retries"`
		RetryDelay  int    `yaml:"retry_delay"`
		SkipVerify  bool   `yaml:"skip_verify"`
	} `yaml:"client"`
}

func LoadServerConfig(filename string) (*ServerConfig, error) {
	config := &ServerConfig{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadClientConfig(filename string) (*ClientConfig, error) {
	config := &ClientConfig{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}