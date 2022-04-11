// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2021

package env_test

import (
	"log"

	"github.com/dtgorski/env"
)

func ExampleUnmarshal() {
	type Config struct {
		MySQL struct {
			Host     string `env:"MYSQL_HOST"`
			Username string `env:"MYSQL_USER"`
			Password string `env:"MYSQL_PASSWORD,file"` // fallback: <(cat $MYSQL_PASSWORD_FILE)
			Database string `env:"MYSQL_DATABASE"`
		}
		Nodes []string `env:"NODES"`
	}

	conf := Config{}

	if err := env.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	println(conf.MySQL.Username, conf.MySQL.Password)
}
