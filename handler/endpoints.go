package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"log"
	"net/http"
	"swpr/config"
	"swpr/generated"
	"swpr/repository"
	"swpr/util"
	"time"
)

func (s *Server) Register(ctx echo.Context) (err error) {
	req := new(generated.RegisterReq)
	if err = ctx.Bind(req); err != nil {
		log.Println("Error ", err)
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
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

	hashPwd, err := s.PasswordUtil.HashPassword(req.Password)
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
	req := new(generated.LoginReq)
	if err = ctx.Bind(req); err != nil {
		log.Println("Error ", err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	//Requirement no.3
	//Check database whether the combination exists (phone & password).
	user, err := s.Repository.UserGetByPhone(ctx.Request().Context(), req.Phone)
	if err != nil {
		log.Println("error ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: []interface{}{
					"Invalid phone or password",
				},
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	//Requirement no.3
	//Capture data login attempt
	var isLoginSuccess = false
	defer func() {
		if user != nil {
			_, _ = s.Repository.LoginAttemptCreate(ctx.Request().Context(), repository.LoginAttemptCreate{
				UserID: user.Id,
				Status: isLoginSuccess,
			})
		}
	}()

	//Requirement no.3
	//Compute param plain password with hash password on user result.
	//Unsuccessful login will return HTTP 400 Bad Requests code.
	isMatch := s.PasswordUtil.CheckPasswordHash(req.Password, user.Password)
	if !isMatch {
		err = errors.New("Invalid Credential")
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	//Requirement no.3
	//JWT with algorithm RS256.
	stringToken, err := util.JwtCreateToken(user.Id, config.Instance.Security.JwtSecKey, config.Instance.Security.JwtTTL)
	if err != nil {
		log.Println("error ", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}
	isLoginSuccess = true

	return ctx.JSON(http.StatusOK, generated.LoginResOk{
		Id:    int(user.Id),
		Token: stringToken,
	})
}

func (s *Server) Profile(ctx echo.Context) error {
	//Requirement no.4
	//JWT as a bearer token in the authorization. Otherwise, return HTTP 403 Forbidden code.
	jwt, err := getJwtFromHeaders(ctx.Request().Header)
	if err != nil {
		log.Println("error getJwtFromHeaders", err)
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	userID, err := util.JwtGetSubjectUserID(jwt, config.Instance.Security.JwtSecKey)
	if err != nil {
		log.Println("error JwtGetSubjectUserID", err)
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	user, err := s.Repository.UserGetById(ctx.Request().Context(), int64(userID))
	if err != nil {
		log.Println("error ", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	//Requirement no.4
	//Upon successful verification, this will return the name of the user and the phone number.
	return ctx.JSON(http.StatusOK, generated.GetProfileResOk{
		Name:  user.FullName,
		Phone: user.Phone,
	})
}

func (s *Server) ProfileUpdate(ctx echo.Context) error {
	//Requirement no.5
	//Accept JWT as bearer token in authorization header.
	//If the request is unauthorised, then return HTTP 403 Forbidden code.
	jwt, err := getJwtFromHeaders(ctx.Request().Header)
	if err != nil {
		log.Println("error getJwtFromHeaders", err)
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}
	userID, err := util.JwtGetSubjectUserID(jwt, config.Instance.Security.JwtSecKey)
	if err != nil {
		log.Println("error JwtGetSubjectUserID", err)
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}
	user, err := s.Repository.UserGetById(ctx.Request().Context(), int64(userID))
	if err != nil {
		log.Println("error ", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	req := new(generated.UpdateProfileReq)
	if err = ctx.Bind(req); err != nil {
		log.Println("Error ", err)
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	var errMsgsValidation []interface{}

	//Requirement no.1 validation
	//Phone format, name format & password format
	if req.Phone != nil {
		isPhoneValid := util.ValidatePhoneFormat(*req.Phone)
		if !isPhoneValid {
			errMsgsValidation = append(errMsgsValidation, "invalid phone number format")
		} else {
			user.Phone = *req.Phone
		}
	}

	if req.FullName != nil {
		isNameValid := util.ValidateFullNameFormat(*req.FullName)
		if !isNameValid {
			errMsgsValidation = append(errMsgsValidation, "invalid fullname. Min.3 Max.60")
		} else {
			user.FullName = *req.FullName
		}
	}

	//Requirement no.2, status bad request, return all error fields, The password should be hashed,
	if len(errMsgsValidation) != 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: errMsgsValidation,
		})
	}

	_, err = s.Repository.UserUpdate(ctx.Request().Context(), repository.UserUpdate{
		Id:        user.Id,
		FullName:  user.FullName,
		Password:  user.Password,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		log.Println("error ", err)

		//Requirement no.5
		//Already exist phone number returns HTTP 409 Conflict.
		var driverErr *pq.Error
		ok := errors.As(err, &driverErr)
		if ok {
			if driverErr.Code == "23505" { //pq: duplicate key value violates unique constraint
				return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
					Message: []interface{}{
						fmt.Sprintf("Error pq code %v", driverErr.Code),
						err.Error(),
						driverErr.Detail,
					},
				})
			}
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: []interface{}{
				err.Error(),
			},
		})
	}

	//Requirement no.5
	return ctx.JSON(http.StatusOK, generated.UpdateProfileResOk{
		Id:      userID,
		Success: true,
	})
}
