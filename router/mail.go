package router

import (
	"lain/controllers"
	"lain/types"
	"lain/utils/auth"
	"lain/utils/urls"
)

func init() {
	urls.SetNamespace("mail")

	urls.Path(types.GET, "/inbox", auth.RequireAuthentication(controllers.Mailbox), "inbox")
	urls.Path(types.GET, "/*", auth.RequireAuthentication(controllers.Mailbox), "folder")
}
