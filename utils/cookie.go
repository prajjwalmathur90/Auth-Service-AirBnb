package utils

import (
	env "AuthInGo/config/env"
	"net/http"
	"time"
)

const (
	// AuthCookieName is the name of the cookie that carries the JWT.
	AuthCookieName = "auth_token"

	// AuthTokenTTL is the lifetime of the JWT and of the cookie that carries it.
	AuthTokenTTL = 24 * time.Hour
)

// SetAuthCookie stores the JWT in an HttpOnly cookie so that JavaScript
// running in the browser can never read it (mitigates XSS token theft).
// Set COOKIE_SECURE=true in production so it is only sent over HTTPS.
func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(AuthTokenTTL),
		MaxAge:   int(AuthTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   env.GetBool("COOKIE_SECURE", false),
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearAuthCookie tells the browser to delete the auth cookie. The attributes
// must match the ones used when the cookie was set, otherwise the browser
// treats it as a different cookie and keeps the original.
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   env.GetBool("COOKIE_SECURE", false),
		SameSite: http.SameSiteLaxMode,
	})
}
