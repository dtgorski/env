# env

Painlessly initialize structs (e.g. Configuration) from process environment variables and files. 

## Installation
```
go get -u github.com/dtgorski/env
```

## Usage
```go
package main

import (
	"log"
	"github.com/dtgorski/env"
)

type Config struct {
	MySQL struct {
		Host     string `env:"MYSQL_HOST"`
		Username string `env:"MYSQL_USER"`
		Password string `env:"MYSQL_PASSWORD,file"` // fallback: <(cat $MYSQL_PASSWORD_FILE)
		Database string `env:"MYSQL_DATABASE"`
		Timeout  int    `env:"MYSQL_TIMEOUT"`
	}
	Nodes []string `env:"NODES"` // slice from comma separated items
}

func main() {
	conf := Config{}

	if err := env.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	println(conf.MySQL.Username, conf.MySQL.Password)
}
```

## Reading values from file
Setting the ``file`` option of a field tag (see Config.MySQL.Password above) will
additionally lookup the environment variable with a ```_FILE``` suffix and load
the contents of the file it points to.

## @dev
Try ```make```:
```
$ make

 make help       Displays this list
 make clean      Removes build/test artifacts
 make test       Runs tests with -race  (pick: ARGS="-run=<Name>")
 make bench      Artificial benchmarks  (pick: ARGS="-bench=<Name>")
 make prof-cpu   Creates CPU profile    (pick: ARGS="-bench=<Name>")
 make prof-mem   Creates memory profile (pick: ARGS="-bench=<Name>")
 make escape     Displays heap escape analysis
 make sniff      Checks format and runs linter (void on success)
 make tidy       Formats source files, cleans go.mod

 Usage: make <TARGET> [ARGS=...]
```

## Disclaimer
The implementation and features of ```env``` follow the [YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it) principle.
There is no claim for completeness or reliability.

## License
[MIT](https://opensource.org/licenses/MIT) - © dtg [at] lengo [dot] org
