package session

import (
	"fmt"
	"lain/config"
	"log"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

var Store *session.Store

func init() {
	storage := postgres.New(postgres.Config{
		Host:       config.Database.Host,
		Port:       config.Database.Port,
		Username:   config.Database.Username,
		Password:   config.Database.Password,
		Database:   config.Database.Name,
		Table:      config.Session.CookieName,
		Reset:      false,
		SSLMode:    config.Database.SSLMode,
		GCInterval: 10 * time.Second,
	})

	Store = session.New(session.Config{
		Storage:        storage,
		Expiration:     config.Session.CookieTimeout,
		KeyLookup:      fmt.Sprintf("cookie:%s", config.Session.CookieName),
		CookieDomain:   config.Session.CookieDomain,
		CookiePath:     config.Session.CookiePath,
		CookieSecure:   config.Session.CookieSecure,
		CookieSameSite: config.Session.CookieSameSite,
		CookieHTTPOnly: true,
	})

	log.Println("session storage initialized")
}
