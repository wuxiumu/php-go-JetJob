package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config 统一配置对象
type Config struct {
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
	Redis struct {
		Addr string `yaml:"addr"`
	} `yaml:"redis"`
	Token         string `yaml:"token"`
	MaxConcurrent int    `yaml:"max_concurrent"`
}

var Cfg Config

// LoadConfig 从yaml加载全局配置
func LoadConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&Cfg); err != nil {
		panic(err)
	}
}
