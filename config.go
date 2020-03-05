package sweb

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Ssl      bool   `toml:"ssl"`
	CertFile string `toml:"cert-file"`
	KeyFile  string `toml:"key-file"`
	Address  string `toml:"address"`
	Root     string `toml:"root"`
	Name     string `toml:"name"`
}

type OauthConfig struct {
	ClientId string `toml:"client-id"`
	ClientSecret string `toml:"client-secret"`
}

func FindConfigPath(name string) string {
	return "conf/" + name + "/" + name + ".toml"
}

func saveConfig(path string, v interface{}) {
	f, err := os.Create(path)

	if err != nil {
		log.Panicln(err)
	}

	e := toml.NewEncoder(f)

	if err := e.Encode(v); err != nil {
		log.Panicln(err)
	}
}

func LoadConfig(path string, v interface{}) bool {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Panicln(err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		saveConfig(path, v)

		return false
	}

	if _, err := toml.DecodeFile(path, v); err != nil {
		log.Panicln(err)
	}

	return true
}
