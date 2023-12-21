package handler

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"log"
	"net/http"
	"strings"
	"swpr/config"
	"swpr/generated"
	"swpr/util"
)

func JwtMiddlewareValidation() *middleware.Options {
	return &middleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			log.Println("ErrorHandler")

			if err.Code == 401 { // hndling error 401 dari AuthenticationFunc
				return c.JSON(http.StatusForbidden, generated.ErrorResponse{
					Message: []interface{}{
						"Forbidden Access",
					},
				})
			}
			return nil
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) (err error) {
				log.Println("AuthenticationFunc")

				authorizationHeader := input.RequestValidationInput.Request.Header["Authorization"]
				if authorizationHeader == nil {
					return echo.NewHTTPError(401, "Invalid authorization headers")
				}

				authHeader := strings.Split(authorizationHeader[0], " ")
				if len(authHeader) < 2 {
					return echo.NewHTTPError(401, "Invalid bearer format")
				}

				jwt := authHeader[1]
				isVerify, err := util.JwtVerify(jwt, config.Instance.Security.JwtSecKey)
				if err != nil {
					log.Println("error JwtVerify", err)
					return echo.NewHTTPError(401, "Invalid token")
				}
				if !isVerify {
					log.Println("error JwtVerify isVerify", isVerify)
					return echo.NewHTTPError(401, "Invalid token")
				}

				return nil
			},
		},
	}
}
