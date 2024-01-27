package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// パッケージ由来のstruct
var DB *gorm.DB

func Connect() {

	var err error

	DB, err = gorm.Open(mysql.Open("user:user_password@tcp(db:3306)/ambassador?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic("Could not connect with tha database!")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Product{}, models.Link{}, models.Order{}, models.OrderItem{})
}
