package services

import (
	"app/dto"
	models "app/models/generated"
	"app/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestAuthServiceSuite struct {
	WithDBSuite
}

var testAuthService AuthService

func (s *TestAuthServiceSuite) SetupTest() {
	s.SetDBCon()

	testAuthService = NewAuthService(DBCon)
}

func (s *TestAuthServiceSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestAuthServiceSuite) TestSignUp() {
	requestParams := dto.SignUpRequest{Name: "test name 1", Email: "test@example.com", Password: "password"}

	result := testAuthService.SignUp(ctx, requestParams)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)

	// NOTE: ユーザが作成されていることを確認
	isExistUser, err := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	if err != nil {
		s.T().Fatalf("failed to create user %v", err)
	}
	assert.True(s.T(), isExistUser)
}

func (s *TestAuthServiceSuite) TestSignUp_ValidationError() {
	requestParams := dto.SignUpRequest{Name: "test name 1", Email: "", Password: "password"}

	result := testAuthService.SignUp(ctx, requestParams)

	assert.NotNil(s.T(), result.Error)
	assert.Equal(s.T(), "validationError", result.ErrorType)

	// NOTE: ユーザが作成されていないことを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ?", "test name 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistUser)
}

func (s *TestAuthServiceSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := dto.SignInRequest{Email: "test@example.com", Password: "password"}

	result := testAuthService.SignIn(ctx, requestParams)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.NotFoundMessage)
	assert.NotNil(s.T(), result.TokenString)
}

func (s *TestAuthServiceSuite) TestSignIn_NotFoundError() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := dto.SignInRequest{Email: "test_1@example.com", Password: "password"}

	result := testAuthService.SignIn(ctx, requestParams)

	assert.Equal(s.T(), "メールアドレスまたはパスワードに該当するユーザが存在しません。", result.NotFoundMessage)
}

func TestAuthService(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestAuthServiceSuite))
}
