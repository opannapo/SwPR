package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"swpr/generated"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostRegister(ctx echo.Context) error {
	var resp generated.RegisterRes
	resp.Success = true
	return ctx.JSON(http.StatusOK, resp)
}
