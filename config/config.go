package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Load(filename string) (Config, error) {
	cfg := Config{}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func LoadFromStr(str string) (Config, error) {
	cfg := Config{}

	if err := json.Unmarshal([]byte(str), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func LoadDefault() (Config, error) {
	return Load(filepath.Dir(os.Args[0]) + "/config/application.json")
}

type WebConfig struct {
	Port    string `json:"port"`
	IsDebug bool   `json:"isDebug"`
}

type DbConfig struct {
	Type               string `json:"type"`
	Dsn                string `json:"dsn"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
}

type CacheConfig struct {
	Expiration int32  `json:"expiration"`
	Prefix     string `json:"prefix"`
	Host       string `json:"host"`
}

type RedisConfig struct {
	Address     string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Wait        bool
}

type Config struct {
	Web   WebConfig              `json:"web"`
	Db    map[string]DbConfig    `json:"db"`
	Cache CacheConfig            `json:"cache"`
	Redis map[string]RedisConfig `json:"redis"`
}
