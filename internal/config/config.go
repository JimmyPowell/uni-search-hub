package config

import (
	"fmt"
	"os"
	"strings"
	"uni-search-hub/pkg/common"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig 存储服务器相关的配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 存储所有数据库连接的配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

// MySQLConfig 存储 MySQL 数据库的配置。
type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

// RedisConfig 存储 Redis 的配置。
type RedisConfig struct {
	Addr          string `mapstructure:"addr"`
	Password      string `mapstructure:"password"`
	SyncFrequency string `mapstructure:"sync_frequency"`
}

// LoadConfig 加载配置，从指定路径读取 YAML 配置文件并解析到 fx 框架中
func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Errorf("配置文件不存在: %s", configPath))
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("无法将配置解析到结构体中: %w", err))
	}

	return config
}

func initEnv(config *Config) {
	common.DebugEnabled = config.Server.Mode == "debug"
}
