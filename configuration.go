package main

import (
	"github.com/BurntSushi/toml"
)

// DBConfiguration struct
type DBConfiguration struct {
	Host     string `toml:"port"`
	Name     string `toml:"name"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
}

// HTTPConfiguration struct
type HTTPConfiguration struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// Config struct
type Config struct {
	DB   *DBConfiguration   `toml:"db"`
	HTTP *HTTPConfiguration `toml:"http"`
}

// ConfigFromFile func
func ConfigFromFile(path string) (*Config, error) {
	config := Config{}
	_, err := toml.DecodeFile(path, &config)
	return &config, err
}
