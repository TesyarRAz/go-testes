package network

import (
	"fmt"

	"github.com/TesyarRAz/testes/domain/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(driverName, host, port, name, user, password string) (*Database, error) {
	var driver gorm.Dialector

	if driverName == "sqlite" {
		driver = sqlite.Open(host)
	} else {
		url := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, name, user, password)

		if driverName == "mysql" {
			driver = mysql.Open(url)
		}
	}

	db, err := gorm.Open(driver, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) AutoMigrate() error {
	return db.Client.AutoMigrate(
		&entity.User{},
	)
}
