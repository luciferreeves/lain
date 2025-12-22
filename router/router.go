package router

import (
	"lain/controllers"
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Static("/static", "./static")

	urls.Attach(router)
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	switch code {
	case fiber.StatusBadRequest:
		return controllers.BadRequest(ctx, err)
	case fiber.StatusForbidden:
		return controllers.Forbidden(ctx, err)
	case fiber.StatusInternalServerError:
		return controllers.InternalServerError(ctx, err)
	case fiber.StatusNotFound:
		return controllers.NotFound(ctx, err)
	case fiber.StatusUnauthorized:
		return controllers.Unauthorized(ctx, err)
	default:
		return controllers.DefaultError(ctx, err)
	}
}
