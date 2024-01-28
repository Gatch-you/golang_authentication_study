package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id           uint     `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email" gorm:"unique"`
	Password     []byte   `json:"-"`
	IsAmbassador bool     `json:"-"`
	Revenure     *float64 `json:"revenure,omitempty" gorm:"-"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

type Admin User

func (admin *Admin) CalculateRevenure(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   admin.Id,
		Complete: true,
	})

	var revenure float64 = 0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenure += orderItem.AdminRevenue
		}
	}

	admin.Revenure = &revenure

}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenure(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   ambassador.Id,
		Complete: true,
	})

	var revenure float64 = 0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenure += orderItem.AmbassadorRevenue
		}
	}

	ambassador.Revenure = &revenure
}
