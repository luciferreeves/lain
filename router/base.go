package router

import (
	"lain/types"
	"lain/utils/shortcuts"
	"lain/utils/urls"
)

func init() {
	urls.SetNamespace("")

	urls.Path(types.GET, "/", shortcuts.RedirectTo("auth.login"), "home")
}
