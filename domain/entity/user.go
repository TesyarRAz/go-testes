package entity

import (
	"html"
	"strings"
	"time"

	"github.com/TesyarRAz/testes/infrastructure/security"
)

type User struct {
	ID        uint64     `gorm:"primarykey;auto_increment" json:"id"`
	FirstName string     `gorm:"not null" json:"first_name"`
	LastName  string     `gorm:"not null" json:"last_name"`
	Username  string     `gorm:"not null;unique" json:"username"`
	Password  string     `gorm:"not null" json:"password"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PublicUser struct {
	ID        uint64 `gorm:"primarykey;auto_increment" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))

	for i, user := range users {
		result[i] = user.PublicUser()
	}

	return result
}

func (u *User) PublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) Prepare() {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}
