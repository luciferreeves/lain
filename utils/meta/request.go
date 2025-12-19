package meta

import (
	"lain/types"

	"github.com/gofiber/fiber/v2"
)

func BuildRequest(context *fiber.Ctx) types.HTTPRequest {
	return types.HTTPRequest{
		Path:        context.Path(),
		Method:      context.Method(),
		Query:       buildQueryParams(context),
		Params:      buildRouteParams(context),
		QueryString: string(context.Request().URI().QueryString()),
		IP:          context.IP(),
		URL:         context.OriginalURL(),
	}
}

func buildQueryParams(context *fiber.Ctx) []types.HTTPQueryParam {
	params := make([]types.HTTPQueryParam, 0)
	args := context.Request().URI().QueryArgs()

	args.VisitAll(transformQueryParam(&params))
	return params
}

func buildRouteParams(context *fiber.Ctx) []types.HTTPQueryParam {
	params := make([]types.HTTPQueryParam, 0)
	for key, value := range context.AllParams() {
		params = append(params, types.HTTPQueryParam{
			Key:   key,
			Value: value,
		})
	}
	return params
}

func transformQueryParam(params *[]types.HTTPQueryParam) func(key, value []byte) {
	return func(key, value []byte) {
		*params = append(*params, types.HTTPQueryParam{
			Key:   string(key),
			Value: string(value),
		})
	}
}
