package main

import (
	"ambassador/src/database"
	"ambassador/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// arrow_origin := os.Getenv("ACCESS_ALLOW_ORIGIN")

	database.Connect()
	database.AutoMigrate()
	database.SetupRedis()
	database.SetupCacheChannel()

	app := fiber.New()

	// どうやらReactアプリケーション間の通信はワイルドカード使えないようなので自身で設定しとく
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	routes.Setup(app)

	app.Listen(":3000")
}
