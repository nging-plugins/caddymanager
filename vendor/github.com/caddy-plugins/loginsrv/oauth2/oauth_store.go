package oauth2

import (
	"net/http"
	"os"

	"github.com/admpub/goth/gothic"
	"github.com/gorilla/sessions"
	"github.com/webx-top/com"
)

func init() {
	if len(os.Getenv("SESSION_SECRET")) > 0 {
		return
	}
	SetCookieStore(com.RandomAlphanumeric(16))
}

func SetCookieStore(secret string) {
	cookieStore := sessions.NewCookieStore([]byte(secret))
	cookieStore.Options.HttpOnly = true
	SetSessionStore(cookieStore)
}

func SetSessionStore(store sessions.Store) {
	gothic.Store = store
}

// Logout invalidates a user session.
func Logout(res http.ResponseWriter, req *http.Request) error {
	return gothic.Logout(res, req)
}
