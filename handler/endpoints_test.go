package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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
	"swpr/repository"
	"testing"
)

type HandlerTestSuite struct {
	suite.Suite
	mockRepository *repository.MockRepositoryInterface
	config         *config.AppConfig
	Server         generated.ServerInterface
	Echo           *echo.Echo
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
	h.mockRepository = repository.NewMockRepositoryInterface(mockCtrl)

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
		Repository: h.mockRepository,
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

			},
			expectedResponse: func(actualRes *http.Response) {
				assert.Equal(h.T(), 200, actualRes.StatusCode)
				resMap := h.parseResponseJson(actualRes)
				assert.NotEmpty(h.T(), resMap["id"])
				assert.NotEmpty(h.T(), resMap["token"])
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
