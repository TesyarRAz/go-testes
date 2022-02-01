package application

import (
	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/TesyarRAz/testes/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

type UserAppInterface interface {
	Save(*entity.User) error
	Find(uint64) (*entity.User, error)
	Get() (entity.Users, error)
	FindUserByUsernameAndPassword(*entity.User) (*entity.User, error)
}

var _ UserAppInterface = &userApp{}

func (u *userApp) Save(user *entity.User) error {
	return u.us.Save(user)
}

func (u *userApp) Find(id uint64) (*entity.User, error) {
	return u.us.Find(id)
}

func (u *userApp) Get() (entity.Users, error) {
	return u.us.Get()
}

func (u *userApp) FindUserByUsernameAndPassword(user *entity.User) (*entity.User, error) {
	return u.us.FindUserByUsernameAndPassword(user)
}
