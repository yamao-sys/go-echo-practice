package services

import (
	"app/dto"
	models "app/models/generated"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx context.Context, requestParams dto.SignUpRequest) *dto.SignUpResponse
	SignIn(ctx context.Context, requestParams dto.SignInRequest) *dto.SignInResponse
	GetAuthUser(ctx echo.Context) (*models.User, error)
	Getuser(ctx context.Context, id int) *models.User
}

type authService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) AuthService {
	return &authService{db}
}

func (as *authService) SignUp(ctx context.Context, requestParams dto.SignUpRequest) *dto.SignUpResponse {
	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(requestParams)
	if validationErrors != nil {
		return &dto.SignUpResponse{User: models.User{}, Error: validationErrors, ErrorType: "validationError"}
	}

	user := models.User{}
	user.Name = requestParams.Name
	user.Email = requestParams.Email
	// NOTE: パスワードをハッシュ化の上、Create処理
	hashedPassword, err := as.encryptPassword(requestParams.Password)
	if err != nil {
		return &dto.SignUpResponse{User: user, Error: err, ErrorType: "internalServerError"}
	}
	user.Password = hashedPassword
	createErr := user.Insert(ctx, as.db, boil.Infer())
	if createErr != nil {
		log.Fatalln(createErr)
	}

	return &dto.SignUpResponse{User: user, Error: nil, ErrorType: ""}
}

func (as *authService) SignIn(ctx context.Context, requestParams dto.SignInRequest) *dto.SignInResponse {
	// NOTE: emailからユーザの取得
	user, err := models.Users(qm.Where("email = ?", requestParams.Email)).One(ctx, as.db)
	if err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "メールアドレスまたはパスワードに該当するユーザが存在しません。", Error: nil}
	}

	// NOTE: パスワードの照合
	if err := as.compareHashPassword(user.Password, requestParams.Password); err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "メールアドレスまたはパスワードに該当するユーザが存在しません。", Error: nil}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	// TODO: JWT_SECRETを環境変数に切り出す
	tokenString, err := token.SignedString([]byte("abcdefghijklmn"))
	if err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "", Error: err}
	}
	return &dto.SignInResponse{TokenString: tokenString, NotFoundMessage: "", Error: nil}
}

func (as *authService) GetAuthUser(ctx echo.Context) (*models.User, error) {
	// NOTE: Cookieからtokenを取得
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		return &models.User{}, err
	}
	// NOTE: tokenに該当するユーザを取得する
	token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("abcdefghijklmn"), nil
	})
	if err != nil {
		return &models.User{}, fmt.Errorf("failt jwt parse")
	}

	var userID int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = int(claims["user_id"].(float64))
	}
	if userID == 0 {
		return &models.User{}, fmt.Errorf("invalid token")
	}
	user, err := models.FindUser(context.Background(), as.db, userID)
	return user, err
}

func (as *authService) Getuser(ctx context.Context, id int) *models.User {
	user, _ := models.FindUser(ctx, as.db, id)
	return user
}

// NOTE: パスワードの文字列をハッシュ化する
func (as *authService) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// NOTE: パスワードの照合
func (as *authService) compareHashPassword(hashedPassword, requestPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword)); err != nil {
		return err
	}
	return nil
}
