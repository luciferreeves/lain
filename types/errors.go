package types

import "github.com/gofiber/fiber/v2"

type TemplateError struct {
	Context      *fiber.Ctx
	PageTitle    string
	ErrorMessage error
	StatusCode   int
}
