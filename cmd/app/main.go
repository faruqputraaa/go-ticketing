package main

import (
	"fmt"

	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/pkg/database"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	_, err = database.InitDatabase(cfg.PostgresConfig)
	checkError(err)
	fmt.Println("database connected")

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
