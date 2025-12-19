package router

import (
	"lain/controllers"
	"lain/types"
	"lain/utils/urls"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(types.GET, "/login", controllers.LoginPage, "login")
}
