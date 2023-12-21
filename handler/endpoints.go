package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"swpr/generated"
	"swpr/repository"
	"swpr/util"
)

func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Register(ctx echo.Context) (err error) {
	req := new(generated.RegisterReq)
	if err = ctx.Bind(req); err != nil {
		log.Println("Error ", err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	var errMsgsValidation []interface{}

	//validation
	isPhoneValid := util.ValidatePhoneForRegister(req.Phone)
	if !isPhoneValid {
		errMsgsValidation = append(errMsgsValidation, "invalid phone number format")
	}
	isNameValid := util.ValidateNameForRegister(req.FullName)
	if !isNameValid {
		errMsgsValidation = append(errMsgsValidation, "invalid fullname. Min.3 Max.60")
	}

	if len(errMsgsValidation) != 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: errMsgsValidation,
		})
	}

	idResult, err := s.Repository.UserCreate(ctx.Request().Context(), repository.UserCreate{
		FullName: req.FullName,
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		log.Println("error ", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, generated.RegisterResOk{
		Success: true,
		Id:      int(idResult),
	})
}
