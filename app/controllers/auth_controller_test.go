package controllers

import (
	models "app/models/generated"
	"app/services"
	"app/test/factories"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	testAuthController AuthController
)

type TestAuthControllerSuite struct {
	WithDBSuite
}

func (s *TestAuthControllerSuite) SetupTest() {
	s.SetDBCon()

	authService := services.NewAuthService(DBCon)

	// NOTE: テスト対象のコントローラを設定
	testAuthController = NewAuthController(authService)
}

func (s *TestAuthControllerSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestAuthControllerSuite) TestSignUp() {
	echoServer := echo.New()
	res := httptest.NewRecorder()
	signUpRequestBody := bytes.NewBufferString("{\"name\":\"test name 1\",\"email\":\"test@example.com\",\"password\":\"password\"}")
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_up", signUpRequestBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echoServer.NewContext(req, res)
	c.SetPath("auth/sign_up")
	testAuthController.SignUp(c)

	assert.Equal(s.T(), 200, res.Code)

	// NOTE: ユーザが作成されていることを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistUser)
}

func (s *TestAuthControllerSuite) TestSignUp_ValidationError() {
	echoServer := echo.New()
	res := httptest.NewRecorder()
	signUpRequestBody := bytes.NewBufferString("{\"name\":\"test name 1\",\"email\":\"\",\"password\":\"password\"}")
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_up", signUpRequestBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echoServer.NewContext(req, res)
	c.SetPath("auth/sign_up")
	testAuthController.SignUp(c)

	assert.Equal(s.T(), 400, res.Code)

	// NOTE: ユーザが作成されていないことを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistUser)
}

func (s *TestAuthControllerSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	echoServer := echo.New()
	res := httptest.NewRecorder()
	signInRequestBody := bytes.NewBufferString("{\"email\":\"test@example.com\",\"password\":\"password\"}")
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_in", signInRequestBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echoServer.NewContext(req, res)
	c.SetPath("auth/sign_in")
	testAuthController.SignIn(c)

	assert.Equal(s.T(), 200, res.Code)
	token = res.Result().Cookies()[0].Value
	assert.NotEmpty(s.T(), token)
}

func (s *TestAuthControllerSuite) TestSignIn_NotFoundError() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	echoServer := echo.New()
	res := httptest.NewRecorder()
	signInRequestBody := bytes.NewBufferString("{\"email\":\"test_1@example.com\",\"password\":\"password\"}")
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_in", signInRequestBody)
	req.Header.Set(echo.HeaderContentType, "application/json")
	c := echoServer.NewContext(req, res)
	c.SetPath("auth/sing_in")
	testAuthController.SignIn(c)

	assert.Equal(s.T(), 404, res.Code)
	assert.Empty(s.T(), res.Result().Cookies())
}

func TestAuthController(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestAuthControllerSuite))
}
