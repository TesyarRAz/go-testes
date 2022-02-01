package service

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type TokenInterface interface {
	CreateToken(userid uint64) (*entity.TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*entity.AccessDetails, error)
}

type TokenService struct {
}

var _ TokenInterface = &TokenService{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (t *TokenService) CreateToken(userid uint64) (*entity.TokenDetails, error) {
	td := &entity.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.TokenUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + strconv.Itoa(int(userid))

	var (
		token []byte
		err   error
	)
	//Creating Access Token
	atClaims := jwt.New()
	atClaims.Set("authorized", true)
	atClaims.Set("access_uuid", td.TokenUuid)
	atClaims.Set("user_id", userid)
	atClaims.Set("exp", td.AtExpires)
	token, err = jwt.Sign(atClaims, jwa.HS256, []byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	td.AccessToken = string(token)

	//Creating Refresh Token
	rtClaims := jwt.New()
	rtClaims.Set("refresh_uuid", td.RefreshUuid)
	rtClaims.Set("user_id", userid)
	rtClaims.Set("exp", td.RtExpires)
	token, err = jwt.Sign(rtClaims, jwa.HS256, []byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = string(token)

	return td, nil
}

func VerifyToken(r *http.Request) (jwt.Token, error) {
	tokenString := ExtractToken(r)
	return jwt.ParseString(tokenString,
		jwt.WithVerify(jwa.HS256, []byte(os.Getenv("ACCESS_SECRET"))),
		jwt.WithValidate(true),
	)
}

//get the token from the request body
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (t *TokenService) ExtractTokenMetadata(r *http.Request) (*entity.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	accessUuid, ok := token.Get("access_uuid")
	if !ok {
		return nil, err
	}
	userIdString, ok := token.Get("user_id")
	if !ok {
		return nil, err
	}
	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", userIdString), 10, 64)
	if err != nil {
		return nil, err
	}
	return &entity.AccessDetails{
		TokenUuid: accessUuid.(string),
		UserId:    userId,
	}, nil
}
