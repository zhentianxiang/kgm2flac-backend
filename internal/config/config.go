package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Addr            string `yaml:"addr" json:"addr"`
	FFmpegBin       string `yaml:"ffmpeg_bin" json:"ffmpeg_bin"`
	MaxFileSize     int64  `yaml:"max_file_size" json:"max_file_size"`
	MaxFiles        int    `yaml:"max_files" json:"max_files"`
	ParseFormMemory int64  `yaml:"parse_form_memory" json:"parse_form_memory"`
}

// 默认配置
func DefaultConfig() *Config {
	return &Config{
		Addr:            ":8080",
		FFmpegBin:       "ffmpeg",
		MaxFileSize:     1 << 30, // 1GB
		MaxFiles:        50,
		ParseFormMemory: 32 << 20, // 32MB
	}
}

// 加载配置
func LoadConfig(configPath, addr, ffmpegBin string) (*Config, error) {
	cfg := DefaultConfig()

	// 如果指定了配置文件，从文件加载
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// 命令行参数优先于配置文件
	if addr != ":8080" {
		cfg.Addr = addr
	}
	if ffmpegBin != "ffmpeg" {
		cfg.FFmpegBin = ffmpegBin
	}

	return cfg, nil
}

// 保存配置示例
func (c *Config) SaveExample() error {
	example := DefaultConfig()
	data, err := yaml.Marshal(example)
	if err != nil {
		return err
	}

	return os.WriteFile("config.example.yaml", data, 0644)
}
