package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/go-redis/redis/v8"
)

type AuthInterface interface {
	CreateAuth(context.Context, uint64, *entity.TokenDetails) error
	FetchAuth(context.Context, string) (uint64, error)
	DeleteRefresh(context.Context, string) error
	DeleteTokens(context.Context, *entity.AccessDetails) error
}

type AuthService struct {
	client *redis.Client
}

var _ AuthInterface = &AuthService{}

func NewAuthService(client *redis.Client) *AuthService {
	return &AuthService{client}
}

//Save token metadata to Redis
func (as *AuthService) CreateAuth(ctx context.Context, userid uint64, td *entity.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := as.client.Set(ctx, td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := as.client.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

//Check the metadata saved
func (as *AuthService) FetchAuth(ctx context.Context, tokenUuid string) (uint64, error) {
	userid, err := as.client.Get(ctx, tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

//Once a user row in the token table
func (as *AuthService) DeleteTokens(ctx context.Context, authD *entity.AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := as.client.Del(ctx, authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := as.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (as *AuthService) DeleteRefresh(ctx context.Context, refreshUuid string) error {
	//delete refresh token
	if deleted, err := as.client.Del(ctx, refreshUuid).Result(); err != nil || deleted == 0 {
		return err
	}

	return nil
}
