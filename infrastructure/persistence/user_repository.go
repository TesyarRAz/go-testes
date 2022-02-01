package persistence

import (
	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/TesyarRAz/testes/domain/repository"
	"github.com/TesyarRAz/testes/infrastructure/security"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// Implementation
var _ repository.UserRepository = &UserRepository{}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Save(u *entity.User) error {
	if err := r.db.Debug().Create(&u).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Find(id uint64) (*entity.User, error) {
	var user entity.User
	if err := r.db.Debug().First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Get() (entity.Users, error) {
	var users entity.Users
	if err := r.db.Debug().Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindUserByUsernameAndPassword(u *entity.User) (*entity.User, error) {
	var user entity.User
	if err := r.db.Debug().First(&user, "username = ?", user.Username).Error; err != nil {
		return nil, err
	}

	if err := security.VerifyPassword(user.Password, u.Password); err != nil {
		return nil, err
	}

	return &user, nil
}
