package shortcuts

import (
	"errors"
	"lain/types"
	"lain/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func BuildErrorMessage(err error, alternateString string) error {
	if err != nil {
		return err
	}

	return errors.New(alternateString)
}

func RenderError(error types.TemplateError) error {
	meta.SetPageTitle(error.Context, error.PageTitle)
	return RenderWithStatus(error.Context, "error", fiber.Map{
		"ErrorTitle":   error.PageTitle,
		"ErrorMessage": error.ErrorMessage.Error(),
	}, error.StatusCode)
}
