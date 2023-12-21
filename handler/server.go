package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"swpr/generated"
	"swpr/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func CreateBearerMiddleware() ([]echo.MiddlewareFunc, error) {
	api, err := generated.GetSwagger()
	if err != nil {
		log.Println("err", err)
	}

	api.Paths.Validate(context.Background())

	va := func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header
			if header["bearer"] == nil {
				return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
					Message: []interface{}{
						"Forbidden Access",
					},
				})
			}
			return nil
		}
	}

	return []echo.MiddlewareFunc{va}, nil
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
	}
}
