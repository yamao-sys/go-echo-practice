package controllers

import (
	"app/db"
	"app/services"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/DATA-DOG/go-txdb"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type WithDBSuite struct {
	suite.Suite
}

var (
	DBCon *sql.DB
	ctx   context.Context
	token string
)

// func (s *WithDBSuite) SetupSuite()                           {} // テストスイート実施前の処理
// func (s *WithDBSuite) TearDownSuite()                        {} // テストスイート終了後の処理
// func (s *WithDBSuite) SetupTest()                            {} // テストケース実施前の処理
// func (s *WithDBSuite) TearDownTest()                         {} // テストケース終了後の処理
// func (s *WithDBSuite) BeforeTest(suiteName, testName string) {} // テストケース実施前の処理
// func (s *WithDBSuite) AfterTest(suiteName, testName string)  {} // テストケース終了後の処理

func init() {
	txdb.Register("txdb-controller", "mysql", db.GetDsn())
	ctx = context.Background()
}

func (s *WithDBSuite) SetDBCon() {
	db, err := sql.Open("txdb-controller", "connect")
	if err != nil {
		s.T().Fatalf("failed to initialize DB: %v", err)
	}
	DBCon = db
}

func (s *WithDBSuite) CloseDB() {
	DBCon.Close()
}

func (s *WithDBSuite) SignIn() {
	authService := services.NewAuthService(DBCon)
	authController := NewAuthController(authService)

	// recorderの初期化
	authRecorder := httptest.NewRecorder()

	// NOTE: リクエストの生成
	f := make(url.Values)
	f.Set("email", "test@example.com")
	f.Set("password", "password")
	req, _ := http.NewRequest(http.MethodPost, "/auth/sign_in", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// echoによるWebサーバの設定
	echoServer := echo.New()
	c := echoServer.NewContext(req, authRecorder)
	c.SetPath("auth/sign_up")

	// NOTE: ログインし、tokenに認証情報を格納
	authController.SignIn(c)
	token = authRecorder.Result().Cookies()[0].Value
}
