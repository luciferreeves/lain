package session

import "github.com/gofiber/fiber/v2"

func CreateSession(context *fiber.Ctx, email string) error {
	sess, err := Store.Get(context)
	if err != nil {
		return err
	}

	sess.Set("email", email)
	return sess.Save()
}

func DestroySession(context *fiber.Ctx) error {
	sess, err := Store.Get(context)
	if err != nil {
		return err
	}

	return sess.Destroy()
}

func GetSessionEmail(context *fiber.Ctx) (string, error) {
	sess, err := Store.Get(context)
	if err != nil {
		return "", err
	}

	email := sess.Get("email")
	if emailStr, ok := email.(string); ok {
		return emailStr, nil
	}

	return "", nil
}
