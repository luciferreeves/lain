package router

import (
	"lain/controllers"
	"lain/types"
	"lain/utils/urls"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(types.GET, "/login", controllers.LoginPage, "login")
	urls.Path(types.GET, "/logout", controllers.Logout, "logout")

	urls.Path(types.POST, "/login", controllers.Login, "login.submit")
}
