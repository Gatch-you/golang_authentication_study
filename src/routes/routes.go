package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")
	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)
	// admin.Get("user", controllers.User)
	// admin.Post("logout", controllers.Logout)

	adminAuthenticated := admin.Use(middlewares.IsAuthentivated)
	adminAuthenticated.Get("user", controllers.User)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Put("user/info", controllers.UpdateInfo)
	adminAuthenticated.Put("user/password", controllers.UpdatePassword)
	adminAuthenticated.Get("ambassadors", controllers.Ambassadors)
	adminAuthenticated.Get("products", controllers.Products)
	adminAuthenticated.Post("products", controllers.CreateProducts)
	adminAuthenticated.Get("products/:id", controllers.GetProducts)
	adminAuthenticated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProducts)
	adminAuthenticated.Get("users/:id/links", controllers.Link)
	adminAuthenticated.Get("orders", controllers.Orders)

	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)
	ambassador.Get("products/frontend", controllers.ProductsFrontend)
	ambassador.Get("products/backend", controllers.ProductsBackend)

	ambassadorAuthentivated := ambassador.Use(middlewares.IsAuthentivated)
	ambassadorAuthentivated.Get("user", controllers.User)
	ambassadorAuthentivated.Post("logout", controllers.Logout)
	ambassadorAuthentivated.Put("user/info", controllers.UpdateInfo)
	ambassadorAuthentivated.Put("user/password", controllers.UpdatePassword)
	ambassadorAuthentivated.Post("links", controllers.CreateLink)
	ambassadorAuthentivated.Get("stats", controllers.Stats)
	ambassadorAuthentivated.Get("rankings", controllers.Rankings)
}
