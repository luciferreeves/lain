package shortcuts

import (
	"lain/session"

	"github.com/gofiber/fiber/v2"
)

const flashKey = "__flash__"

func Flash(ctx *fiber.Ctx, data any) error {
	normalized, err := normalizeBind(data)
	if err != nil {
		return err
	}

	return session.Set(ctx, flashKey, normalized)
}

func ConsumeFlash(ctx *fiber.Ctx) (any, error) {
	value, err := session.Get(ctx, flashKey)
	if err != nil || value == nil {
		return nil, err
	}

	if err := session.Delete(ctx, flashKey); err != nil {
		return nil, err
	}

	return value, nil
}
