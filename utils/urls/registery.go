package urls

import (
	"sync"

	"lain/types"

	"github.com/gofiber/fiber/v2"
)

type registeredRoute struct {
	method    types.HTTPMethod
	path      string
	handler   fiber.Handler
	namespace string
	name      string
	fullPath  string
}

type routeRegistry struct {
	mutex            sync.Mutex
	currentNamespace string
	routes           map[string]registeredRoute
}

var registry = &routeRegistry{
	routes: make(map[string]registeredRoute),
}
