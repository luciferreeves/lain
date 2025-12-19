package router

import (
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Static("/static", "./static")

	urls.Attach(router)
}
