package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Ethereum struct {
		Url string
	}
	Log struct {
		Level   string
		Fpath   string
		Msize   int
		Mage    int
		Mbackup int
	}
	ContractAddress struct {
		Ca string
	}
	DB struct {
		Host       string
		Database   string
		Collection string
	}
}

func NewConfig(fpath string) *Config {
	c := new(Config)

	if file, err := os.Open(fpath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		//toml 파일 디코딩
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return c
		}
	}
}
