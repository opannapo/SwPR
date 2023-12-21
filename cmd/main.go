package main

import (
	"fmt"
	"swpr/config"
	"swpr/generated"
	"swpr/handler"
	"swpr/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	err := config.InitConfigInstance()
	if err != nil {
		panic(err)
		return
	}
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1111"))
}

func newServer() *handler.Server {
	cfg := config.Instance
	dbDsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}

	return handler.NewServer(opts)
}
