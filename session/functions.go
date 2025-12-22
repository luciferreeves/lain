package session

import "github.com/gofiber/fiber/v2"

const emailKey = "email"

func CreateSession(ctx *fiber.Ctx, email string) error {
	return Set(ctx, emailKey, email)
}

func DestroySession(ctx *fiber.Ctx) error {
	return Delete(ctx, emailKey)
}

func GetSessionEmail(ctx *fiber.Ctx) (string, error) {
	value, err := Get(ctx, emailKey)
	if err != nil || value == nil {
		return "", err
	}

	email, ok := value.(string)
	if !ok {
		return "", nil
	}

	return email, nil
}
