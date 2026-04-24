package cookies

import (
	"akastra-access/internal/app/config"
	"net/http"
	"strings"
)

// helper to determine if running in debug
func isDebugMode() bool {
	return strings.ToLower(config.GetEnv("DEBUG", "false")) == "true"
}

// SetCookie sets a secure cross-subdomain cookie
func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	if maxAge <= 0 {
		maxAge = 86400 // default 1 day
	}

	debug := isDebugMode()

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
	}

	if debug {
		// Local development
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
		// Do not set Domain for localhost
	} else {
		// Production
		cookieDomain := config.GetEnv("COOKIE_DOMAIN", ".akastra.id")
		cookie.Domain = cookieDomain
		cookie.Secure = true 
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, cookie)
}

// GetCookie retrieves a cookie value
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteCookie removes a cookie properly
func DeleteCookie(w http.ResponseWriter, name string) {
	debug := isDebugMode()

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	if debug {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	} else {
		cookieDomain := config.GetEnv("COOKIE_DOMAIN", ".akastra.id")
		cookie.Domain = cookieDomain
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, cookie)
}