package controllers

import (
	"lain/types"
	"lain/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(context *fiber.Ctx, err error) error {
	return shortcuts.RenderError(types.TemplateError{
		Context:      context,
		PageTitle:    "400 – Bad Request",
		ErrorMessage: shortcuts.BuildErrorMessage(err, "The request could not be understood by the server."),
		StatusCode:   fiber.StatusBadRequest,
	})
}

func Forbidden(context *fiber.Ctx, err error) error {
	return shortcuts.RenderError(types.TemplateError{
		Context:      context,
		PageTitle:    "403 – Forbidden",
		ErrorMessage: shortcuts.BuildErrorMessage(err, "You do not have permission to access this resource."),
		StatusCode:   fiber.StatusForbidden,
	})
}

func InternalServerError(context *fiber.Ctx, err error) error {
	return shortcuts.RenderError(types.TemplateError{
		Context:      context,
		PageTitle:    "500 – Internal Server Error",
		ErrorMessage: shortcuts.BuildErrorMessage(err, "An unexpected error occurred on the server."),
		StatusCode:   fiber.StatusInternalServerError,
	})
}

func NotFound(context *fiber.Ctx, err error) error {
	return shortcuts.RenderError(types.TemplateError{
		Context:      context,
		PageTitle:    "404 – Not Found",
		ErrorMessage: shortcuts.BuildErrorMessage(err, "The page you are looking for does not exist."),
		StatusCode:   fiber.StatusNotFound,
	})
}

func Unauthorized(context *fiber.Ctx, err error) error {
	return shortcuts.RenderError(types.TemplateError{
		Context:      context,
		PageTitle:    "401 – Unauthorized",
		ErrorMessage: shortcuts.BuildErrorMessage(err, "You must be logged in to access this resource."),
		StatusCode:   fiber.StatusUnauthorized,
	})
}
