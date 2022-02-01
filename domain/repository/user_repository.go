package repository

import "github.com/TesyarRAz/testes/domain/entity"

type UserRepository interface {
	Save(*entity.User) error
	Find(uint64) (*entity.User, error)
	Get() (entity.Users, error)
	FindUserByUsernameAndPassword(*entity.User) (*entity.User, error)
}
