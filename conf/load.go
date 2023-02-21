package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// 如何把配置映射成config对象

// 从toml格式配置文件加载配置
func LoadConfigFromToml(filePath string) error {
	config = NewDefaultConfig()
	// 读取toml格式的配置
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		return fmt.Errorf("load config from file error, path:%s, %s", filePath, err)
	}
	return nil
}

// 加载环境变量配置
func LoadConfigFromEnv() error {
	config = NewDefaultConfig()
	return env.Parse(config)
}
