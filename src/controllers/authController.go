package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// curl http://localhost:8000/api/admin/register/ -X POST -H "Content-Type: application/json" -d '{"first_name": "a", "last_name": "a", "email": "a@.com", "password": "a", "password_confirm": "a"}'
// todo:ここにメールアドレスがユニークである前提で、登録されているメールアドレスがある場合のエラーハンドリングが必要。
func Register(c *fiber.Ctx) error {
	// {"strint": "string", ...}
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// パスワードの再確認
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password do no match!",
		})
	}

	// パスワードをハッシュ化する -> userモデルにメソッドとして移植
	// password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	// if err != nil {
	// 	return c.JSON(fiber.Map{
	// 		"message": "missisng recieve password in hashed process.",
	// 	})
	// }

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: false,
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ? ", data["email"]).First(&user)

	fmt.Println(user.Password)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// user構造体に対してユーザーのインプットとORM経由でDBからの取得したbyte配列の一致を確認しているmodel/userのメソッドとして移植
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	fmt.Println(payload.Subject, payload.ExpiresAt)

	// user.Idに依存したトークンを生成
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON((fiber.Map{
			"message": "Invalid Credentials",
		}))
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// フロントから送られてきたcookieの情報をもとにそのユーザーのプロファイルを取得する。GETメソッド "認証機能"
// curl --location 'localhost:8000/api/admin/user/' \
// --header 'Cookie: jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDUyOTY1NTIsInN1YiI6IjIifQ.VvhMnDsLVXbMK5GYACOcVdst076S61dKcCZDfaEw12E'
func User(c *fiber.Ctx) error {
	// // middlewaresに移動
	// cookie := c.Cookies("jwt")

	// token, _ := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("secret"), nil
	// })

	// // if err != nil || !token.Valid {
	// // 	c.Status(fiber.StatusUnauthorized)
	// // 	return c.JSON(fiber.Map{
	// // 		"message": "unauthntivated",
	// // 	})
	// // }

	// payload := token.Claims.(*jwt.StandardClaims)

	id, _ := middlewares.GetUser(c)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUser(c)

	user := models.User{
		Id:        id,
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	id, _ := middlewares.GetUser(c)

	user := models.User{
		Id: id,
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	fmt.Printf("Changing password is successed!")
	return c.JSON(user)
}
