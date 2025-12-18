package router

import "github.com/gofiber/fiber/v2"

func Initialize(router *fiber.App) {
	router.Static("/static", "./static")

	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Lain Mail - Present day, present time")
	})
}
