package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Server struct {
		Mode string
		Port string
	}
	Log struct {
		Level   string
		Fpath   string
		Msize   int
		Mage    int
		Mbackup int
	}
	DB struct {
		Host         string
		Database     string
		UsersDB      string
		UserInfoColl string
	}
	Client struct {
		UrlWemix           string
		UrlEth             string
		UrlKlay            string
		WemixTokenAddress  map[string]string
		EthTokenAddress    map[string]string
		KlaytnTokenAddress map[string]string
	}

	Oauth2 map[string]map[string]interface{}
	Wallet map[string]map[string]interface{}
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
			//fmt.Println(c)
			return c
		}
	}
}
