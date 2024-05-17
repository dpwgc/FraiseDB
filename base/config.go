package base

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigModel struct {
	Server NodeConfig  `yaml:"server" json:"server"`
	Join   NodeConfig  `yaml:"join" json:"join"`
	Store  StoreConfig `yaml:"store" json:"store"`
}

type NodeConfig struct {
	Addr     string `yaml:"addr" json:"addr"`
	TcpPort  int    `yaml:"tcp-port" json:"tcpPort"`
	HttpPort int    `yaml:"http-port" json:"httpPort"`
}

type StoreConfig struct {
	Path string `yaml:"path" json:"path"`
}

// InitConfig 加载配置
func InitConfig() error {
	localConfigBytes := loadConfigFile("config.yaml")
	return yaml.Unmarshal(localConfigBytes, &config)
}

func Config() ConfigModel {
	return config
}

// 读取本地配置文件
func loadConfigFile(path string) []byte {
	//加载本地配置
	configBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return configBytes
}
