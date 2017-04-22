package config

import "github.com/BurntSushi/toml"

// Config - Global config of service
type Config struct {
	Address string `json:"address"`
}

// Parse a toml file
func Parse(path string) Config {
	var c Config

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		panic(err)
	}

	return c
}
