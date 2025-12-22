package session

import "github.com/gofiber/fiber/v2"

func Set(ctx *fiber.Ctx, key string, value any) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set(key, value)
	return sess.Save()
}

func Get(ctx *fiber.Ctx, key string) (any, error) {
	sess, err := Store.Get(ctx)
	if err != nil {
		return nil, err
	}

	return sess.Get(key), nil
}

func Delete(ctx *fiber.Ctx, key string) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		return err
	}

	sess.Delete(key)
	return sess.Save()
}
