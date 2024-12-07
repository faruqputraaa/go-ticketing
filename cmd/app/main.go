package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/builder"
	"github.com/faruqputraaa/go-ticket/pkg/database"
	"github.com/faruqputraaa/go-ticket/pkg/server"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	db, err := database.InitDatabase(cfg.PostgresConfig)
	checkError(err)
	fmt.Println(cfg)

	publicRoute := builder.BuildPublicRoute(cfg, db)
	privateRoute := builder.BuildPrivateRoute(cfg, db)

	srv := server.NewServer(cfg, publicRoute, privateRoute)

	runServer(srv, cfg.PORT)
	waitForShutdown(srv)
}

func runServer(srv *server.Server, port string) {
	go func() {
		err := srv.Start(fmt.Sprintf(":%s", port))
		log.Fatal(err)
	}()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func waitForShutdown(srv *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Mematikan server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Gagal mematikan server: %v", err)
	}
	log.Println("Server berhasil dimatikan.")
}
