package urls

import (
	"lain/types"
	"log"

	"github.com/gofiber/fiber/v2"
)

var methodBinders = map[types.HTTPMethod]func(fiber.Router, string, fiber.Handler) fiber.Router{
	types.GET:     func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Get(path, h) },
	types.POST:    func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Post(path, h) },
	types.PUT:     func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Put(path, h) },
	types.PATCH:   func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Patch(path, h) },
	types.DELETE:  func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Delete(path, h) },
	types.OPTIONS: func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Options(path, h) },
	types.HEAD:    func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Head(path, h) },
}

func Attach(app *fiber.App) {
	namespaceGroups := make(map[string]fiber.Router)

	for fullName, route := range registry.routes {
		group, exists := namespaceGroups[route.namespace]
		if !exists {
			group = app.Group("/" + route.namespace)
			namespaceGroups[route.namespace] = group
		}

		binder, ok := methodBinders[route.method]
		if !ok {
			log.Fatalf("%s", "unsupported HTTP method: "+string(route.method))
		}

		fiberRoute := binder(group, route.path, route.handler)
		fiberRoute.Name(fullName)
	}
}
