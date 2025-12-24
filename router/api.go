package router

import (
	"lain/controllers"
	"lain/types"
	"lain/utils/auth"
	"lain/utils/urls"
)

func init() {
	urls.SetNamespace("api")

	urls.Path(types.GET, "/mail/email/:id", auth.RequireAuthentication(controllers.GetEmailAPI), "get_email")
	urls.Path(types.POST, "/mail/email/:id/flag", auth.RequireAuthentication(controllers.ToggleFlagAPI), "toggle_flag")
	urls.Path(types.POST, "/mail/email/:id/read", auth.RequireAuthentication(controllers.MarkEmailAsReadAPI), "mark_read")
}
