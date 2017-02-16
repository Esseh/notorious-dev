// A set of functions for making and maintaining cookies.
package COOKIE

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Deletes a cookie held in the current session by name.
func Delete(res http.ResponseWriter, name string) {
	http.SetCookie(res, &http.Cookie{
		Name:   name,
		MaxAge: -1,
		Path:   "/",
	})
}

// Initializes a cookie into the current session.
func Make(res http.ResponseWriter, name, value string) error {
	mac, err := createHmac(value)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     name,
		Value:    value + "." + base64.RawURLEncoding.EncodeToString(mac),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   sessionTime,
	}
	http.SetCookie(res, c)
	return nil
}

// Retrieves the value located inside of a cookie.
func GetValue(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	val, mac := splitMac(cookie.Value)
	if good := checkMac(val, mac); !good {
		return "", ERROR_NotMatchingHMac
	}
	return val, nil
}
