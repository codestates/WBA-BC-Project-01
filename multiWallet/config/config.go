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
		Host                     string //DB연결 정보
		DaemonDatabase           string //데몬 데이터베이스
		EthColl                  string //이더리움 컬렉션
		WemixColl                string //위믹스 컬렉션
		KlaytnColl               string //클레이튼 컬렉션
		MultiWalletDatabase      string //멀티월렛 데이터베이스
		UserInfoColl             string //유저 컬렉션
		WalletInfoColl           string //월렛 컬렉션
		MultiSigWalletCollection string //다중서명지갑 컬렉션
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
