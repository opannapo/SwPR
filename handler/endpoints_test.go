package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"swpr/config"
	"swpr/generated"
	repository_mock "swpr/mock/repository"
	util_mock "swpr/mock/util"
	"swpr/repository"
	"testing"
)

type HandlerTestSuite struct {
	suite.Suite
	mockRepository   *repository_mock.MockRepositoryInterface
	mockPasswordUtil *util_mock.MockPasswordInterface
	config           *config.AppConfig
	Server           generated.ServerInterface
	Echo             *echo.Echo
}

func (h *HandlerTestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(h.T())
	defer mockCtrl.Finish()

	os.Setenv("APP_DB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "2345")
	os.Setenv("APP_DB_DB", "swpr")
	os.Setenv("APP_DB_USERNAME", "opannapo")
	os.Setenv("APP_DB_PASSWORD", "opannapo")
	os.Setenv("APP_SEC_JWTKEY", "Test123TestKeyJwt")
	os.Setenv("APP_SEC_JWT_TTL", "1h")
	err := config.InitConfigInstance()
	if err != nil {
		panic(err)
	}
	h.config = config.Instance

	//Repo
	h.mockRepository = repository_mock.NewMockRepositoryInterface(mockCtrl)

	//util
	h.mockPasswordUtil = util_mock.NewMockPasswordInterface(mockCtrl)

	h.Echo = echo.New()
	h.Server = h.newServer()

	swagger, err := generated.GetSwagger()
	h.Assert().Nil(err)

	swagger.Servers = nil
	h.Echo.Use(middleware.OapiRequestValidator(swagger))
	h.Echo.Use(echomiddleware.Logger())
	generated.RegisterHandlers(h.Echo, h.Server)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (h *HandlerTestSuite) parseResponseJson(res *http.Response) map[string]interface{} {
	var response map[string]interface{}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return response
}

func (h *HandlerTestSuite) printResult(caseName interface{}, expected interface{}, result interface{}) {
	const colorReset = "\033[0m"
	const colorBlue = "\033[34m"
	const colorGreen = "\033[32m"

	fmt.Printf(colorGreen+"%v Expected "+colorReset+" %+v \n"+
		colorBlue+"COMPARE TO"+colorReset+
		" Actual %+v \n\n\n", caseName, expected, result)
}

func (h *HandlerTestSuite) newServer() *Server {
	opts := NewServerOptions{
		Repository:   h.mockRepository,
		PasswordUtil: h.mockPasswordUtil,
	}
	return NewServer(opts)
}

func (h *HandlerTestSuite) TestServer_Login() {
	type testCase struct {
		caseName         string
		expectedLogic    func(ctx echo.Context, c testCase)
		expectedResponse func(actualRes *http.Response)
		payloadJson      string
	}
	cases := []testCase{
		{
			caseName:    "Success 200",
			payloadJson: `{"password": "111111","phone": "+628561234532"}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetByPhone(ctx.Request().Context(), "+628561234532").
					Return(&userRresult, nil)

				h.mockRepository.EXPECT().
					LoginAttemptCreate(ctx.Request().Context(), repository.LoginAttemptCreate{
						UserID: 1,
						Status: true,
					}).
					Return(int64(1), nil)

				h.mockPasswordUtil.EXPECT().CheckPasswordHash("111111", userRresult.Password).Return(true)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 200, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["id"])
				assert.NotEmpty(h.T(), resMap["token"])
			},
		},
		{
			caseName:    "Error Invalid Password",
			payloadJson: `{"password": "222222","phone": "+628561234532"}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetByPhone(ctx.Request().Context(), "+628561234532").
					Return(&userRresult, nil)

				h.mockRepository.EXPECT().
					LoginAttemptCreate(ctx.Request().Context(), repository.LoginAttemptCreate{
						UserID: 1,
						Status: false,
					}).
					Return(int64(1), nil)

				h.mockPasswordUtil.EXPECT().CheckPasswordHash("222222", userRresult.Password).Return(false)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 400, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "Invalid Credential")
			},
		},
		{
			caseName:    "Error User Not Found",
			payloadJson: `{"password": "222222","phone": "+628561234532"}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				h.mockRepository.EXPECT().
					UserGetByPhone(ctx.Request().Context(), "+628561234532").
					Return(nil, sql.ErrNoRows)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 400, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "Invalid phone or password")
			},
		},
	}

	for _, c := range cases {
		h.Run(c.caseName, func() {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(c.payloadJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			eCtx := h.Echo.NewContext(req, rec)

			c.expectedLogic(eCtx, c)
			_ = h.Server.Login(eCtx)

			//validate detail response
			c.expectedResponse(rec.Result())
		})
	}
}

func (h *HandlerTestSuite) TestServer_Register() {
	type testCase struct {
		caseName         string
		expectedLogic    func(ctx echo.Context, c testCase)
		expectedResponse func(actualRes *http.Response)
		payloadJson      string
	}
	cases := []testCase{
		{
			caseName: "Success 200",
			payloadJson: `{
				"full_name": "opan2",
				"password": "111111",
				"phone": "+628561234511111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				h.mockPasswordUtil.EXPECT().
					HashPassword("111111").
					Return("$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO", nil)

				h.mockRepository.EXPECT().
					UserCreate(ctx.Request().Context(), repository.UserCreate{
						FullName: "opan2",
						Password: "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
						Phone:    "+628561234511111",
					}).
					Return(int64(1), nil)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 200, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["id"])
				assert.NotEmpty(h.T(), resMap["success"])
				assert.Equal(h.T(), resMap["success"], true)
			},
		},
		{
			caseName: "Error phone format : too short",
			payloadJson: `{
				"full_name": "opan2",
				"password": "111111",
				"phone": "+6285612345"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 400, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "invalid phone number format")
			},
		},
		{
			caseName: "Error name format : too short",
			payloadJson: `{
				"full_name": "op",
				"password": "111111",
				"phone": "+62856123455678"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 400, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "invalid fullname. Min.3 Max.60")
			},
		},
		{
			caseName: "Error password format : too short",
			payloadJson: `{
				"full_name": "opan",
				"password": "111",
				"phone": "+62856123455678"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 400, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "invalid password. Min.6 Max.64")
			},
		},
	}

	for _, c := range cases {
		h.Run(c.caseName, func() {
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(c.payloadJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			eCtx := h.Echo.NewContext(req, rec)

			c.expectedLogic(eCtx, c)
			_ = h.Server.Register(eCtx)

			//validate detail response
			c.expectedResponse(rec.Result())
		})
	}
}

func (h *HandlerTestSuite) TestServer_Update() {
	type testCase struct {
		caseName         string
		expectedLogic    func(ctx echo.Context, c testCase)
		expectedResponse func(actualRes *http.Response)
		payloadJson      string
		token            string
	}
	cases := []testCase{
		{
			caseName: "Success 200",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			payloadJson: `{
				"full_name": "opannapo-edit-2",
				"phone": "+628561234511111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, nil)

				h.mockRepository.EXPECT().
					UserUpdate(ctx.Request().Context(), gomock.Any()).
					Return(int64(1), nil)

			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 200, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["id"])
				assert.NotEmpty(h.T(), resMap["success"])
				assert.Equal(h.T(), resMap["success"], true)
			},
		},
		{
			caseName: "Error 409 phone number already exists",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			payloadJson: `{
				"full_name": "opannapo-edit-2",
				"phone": "+628561234511111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, nil)

				driverErr := pq.Error{
					Code: "23505",
				}
				h.mockRepository.EXPECT().
					UserUpdate(ctx.Request().Context(), gomock.Any()).
					Return(int64(1), &driverErr)

			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusConflict, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), errMsg, "Error pq code 23505")
			},
		},
		{
			caseName: "Error invalid token",
			token:    "bearer",
			payloadJson: `{
				"full_name": "opannapo-edit-2",
				"phone": "+628561234511111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusForbidden, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), "Invalid bearer format", errMsg)
			},
		},
		{
			caseName: "Error invalid phone mnumber",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			payloadJson: `{
				"full_name": "opannapo-edit-2",
				"phone": "+628561234511111000000"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, nil)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusBadRequest, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), "invalid phone number format", errMsg)
			},
		},
		{
			caseName: "Error invalid name format",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			payloadJson: `{
				"full_name": "op",
				"phone": "+62856123451111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, nil)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusBadRequest, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), "invalid fullname. Min.3 Max.60", errMsg)
			},
		},
		{
			caseName: "Error db query get user",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			payloadJson: `{
				"full_name": "opan",
				"phone": "+62856123451111"
			}`,
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opan",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+62856123451111",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, sql.ErrConnDone)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusInternalServerError, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), sql.ErrConnDone.Error(), errMsg)
			},
		},
	}

	for _, c := range cases {
		h.Run(c.caseName, func() {
			req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(c.payloadJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, c.token)

			rec := httptest.NewRecorder()
			eCtx := h.Echo.NewContext(req, rec)

			c.expectedLogic(eCtx, c)
			_ = h.Server.ProfileUpdate(eCtx)

			//validate detail response
			c.expectedResponse(rec.Result())
		})
	}
}

func (h *HandlerTestSuite) TestServer_GetProfile() {
	type testCase struct {
		caseName         string
		expectedLogic    func(ctx echo.Context, c testCase)
		expectedResponse func(actualRes *http.Response)
		token            string
	}
	cases := []testCase{
		{
			caseName: "Success 200",
			token:    "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyNzQ0NjUsImlhdCI6MTcwMzE4ODA2NSwic3ViIjoiMSJ9.jUQOUOWwixgvhtbq5AnVbsSTkE7F74uJqn_fvopb95g",
			expectedLogic: func(ctx echo.Context, c testCase) {
				userRresult := repository.UserGet{
					Id:        1,
					FullName:  "opannapo",
					Password:  "$2a$10$sxVsc/YxnHyQvFXeV1L2YuazNV8yEyLOF524o1AlFbSy6wjJx9rkO",
					Phone:     "+628561234532",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				h.mockRepository.EXPECT().
					UserGetById(ctx.Request().Context(), int64(1)).
					Return(&userRresult, nil)
			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 200, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["name"])
				assert.NotEmpty(h.T(), resMap["phone"])
				assert.Equal(h.T(), "opannapo", resMap["name"])
				assert.Equal(h.T(), "+628561234532", resMap["phone"])
			},
		},
		{
			caseName:      "Error invalid token",
			token:         "bearer",
			expectedLogic: func(ctx echo.Context, c testCase) {},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), http.StatusForbidden, actualRes.StatusCode)

				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["message"])

				_, ok := resMap["message"].([]interface{})
				assert.Equal(h.T(), ok, true)

				errMsg := resMap["message"].([]interface{})[0]
				assert.Equal(h.T(), "Invalid bearer format", errMsg)
			},
		},
	}

	for _, c := range cases {
		h.Run(c.caseName, func() {
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			req.Header.Set(echo.HeaderAuthorization, c.token)

			rec := httptest.NewRecorder()
			eCtx := h.Echo.NewContext(req, rec)

			c.expectedLogic(eCtx, c)
			_ = h.Server.Profile(eCtx)

			//validate detail response
			c.expectedResponse(rec.Result())
		})
	}
}
