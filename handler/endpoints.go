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

	//Requirement no.1 validation
	//Phone format, name format & password format
	isPhoneValid := util.ValidatePhoneFormat(req.Phone)
	if !isPhoneValid {
		errMsgsValidation = append(errMsgsValidation, "invalid phone number format")
	}
	isNameValid := util.ValidateFullNameFormat(req.FullName)
	if !isNameValid {
		errMsgsValidation = append(errMsgsValidation, "invalid fullname. Min.3 Max.60")
	}
	isPasswordValid := util.ValidatePasswordFormat(req.Password)
	if !isPasswordValid {
		errMsgsValidation = append(errMsgsValidation, "invalid password. Min.6 Max.64")
	}

	//Requirement no.2, status bad request, return all error fields, The password should be hashed,
	if len(errMsgsValidation) != 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: errMsgsValidation,
		})
	}
	hashPwd, err := util.HashPassword(req.Password)
	idResult, err := s.Repository.UserCreate(ctx.Request().Context(), repository.UserCreate{
		FullName: req.FullName,
		Password: hashPwd,
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

func (s *Server) Login(ctx echo.Context) (err error) {
	req := new(generated.RegisterReq)
	if err = ctx.Bind(req); err != nil {
		log.Println("Error ", err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	//Requirement no.3
	//Check database whether the combination exists (phone & password).
	user, err := s.Repository.UserGetByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		log.Println("error ", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResOk{
		Id: int(user.Id),
	})
}
