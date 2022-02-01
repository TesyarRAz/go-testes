package service

import (
	"fmt"
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
	CreateToken(uint64) (*entity.TokenDetails, error)

	ExtractAccessTokenMetadata(map[string]string) (*entity.AccessDetails, error)
	ExtractRefreshTokenMetadata(map[string]string) (*entity.RefreshDetails, error)
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

func (t *TokenService) ExtractAccessTokenMetadata(r map[string]string) (*entity.AccessDetails, error) {
	var (
		token        jwt.Token
		accessUuid   interface{}
		userIdString interface{}
		userId       uint64

		err error
		ok  bool
	)

	if token, err = VerifyAccessToken(r); err != nil {
		return nil, err
	}
	if accessUuid, ok = token.Get("access_uuid"); !ok {
		return nil, err
	}
	if userIdString, ok = token.Get("user_id"); !ok {
		return nil, err
	}
	if userId, err = strconv.ParseUint(fmt.Sprintf("%.f", userIdString), 10, 64); err != nil {
		return nil, err
	}

	return &entity.AccessDetails{
		TokenUuid: accessUuid.(string),
		UserId:    userId,
	}, nil
}

func (t *TokenService) ExtractRefreshTokenMetadata(r map[string]string) (*entity.RefreshDetails, error) {
	var (
		token        jwt.Token
		refreshUuid  interface{}
		userIdString interface{}
		userId       uint64

		err error
		ok  bool
	)

	if token, err = VerifyAccessToken(r); err != nil {
		return nil, err
	}

	if refreshUuid, ok = token.Get("refresh_uuid"); !ok {
		return nil, err
	}
	if userIdString, ok = token.Get("user_id"); !ok {
		return nil, err
	}
	if userId, err = strconv.ParseUint(fmt.Sprintf("%.f", userIdString), 10, 64); err != nil {
		return nil, err
	}
	return &entity.RefreshDetails{
		TokenUuid: refreshUuid.(string),
		UserId:    userId,
	}, nil
}

func VerifyAccessToken(r map[string]string) (jwt.Token, error) {
	tokenString := ExtractToken(r)
	return jwt.ParseString(tokenString,
		jwt.WithVerify(jwa.HS256, []byte(os.Getenv("ACCESS_SECRET"))),
		jwt.WithValidate(true),
	)
}

func VerifyRefreshToken(r map[string]string) (jwt.Token, error) {
	tokenString := r["refresh_token"]
	return jwt.ParseString(tokenString,
		jwt.WithVerify(jwa.HS256, []byte(os.Getenv("REFRESH_SECRET"))),
		jwt.WithValidate(true),
	)
}

//get the token from the request body
func ExtractToken(r map[string]string) string {
	bearToken := string(r["Authorization"])
	if strArr := strings.Split(bearToken, " "); len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
