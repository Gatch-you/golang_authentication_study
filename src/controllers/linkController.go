package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Link(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var links []models.Link

	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(links)
}
