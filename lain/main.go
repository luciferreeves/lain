package main

import (
	"fmt"
	"lain/config"
	"lain/middleware"
	"lain/processors"
	"lain/router"
	"lain/tags"
	"lain/utils/env"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django/v3"
)

func main() {
	if config.Server.AppSecret == env.Defaults(&config.Server).AppSecret {
		log.Println("Warning: AppSecret is set to a default value which is not secure. Please set a strong random secret in your APP_SECRET environment variable or .env file.")
	}

	tags.Initialize()
	engine := django.New("./templates", ".django")
	engine.Reload(config.Server.DevMode)
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}, // Will be extracted to a separate file later
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{
		CrossOriginEmbedderPolicy: "unsafe-none",
	}))
	app.Use(cors.New())

	processors.Initialize(app)
	router.Initialize(app)
	middleware.Initialize(app)

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting server at %s\n", address)
	log.Fatal(app.Listen(address))
}
