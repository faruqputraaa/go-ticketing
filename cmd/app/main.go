package main

import (
	"fmt"

	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/pkg/database"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	_, err = database.InitDatabase(cfg.PostgresConfig)
	checkError(err)
	fmt.Println(cfg)

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
