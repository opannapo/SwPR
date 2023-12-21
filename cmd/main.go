package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	"log"
	"os"
	"swpr/config"
	"swpr/generated"
	"swpr/handler"
	"swpr/repository"
)

func main() {
	err := config.InitConfigInstance()
	if err != nil {
		panic(err)
		return
	}
	e := echo.New()

	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Println("Error loading swagger spec", err)
		os.Exit(1)
	}
	e.Use(echomiddleware.Logger())
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, handler.JwtMiddlewareValidation()))

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
