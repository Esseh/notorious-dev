package OAUTH

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Esseh/retrievable"
)

// Logs the user in with an OAuth id.
func OAuthLogin(req *http.Request, res http.ResponseWriter, id, first, last, redirect string) {
	err := LoginFromOauth(res, req, id)
	if err == ERROR_NoUser {
		RegisterFromOauth(res, req, id, first, last)
	}
	redirect = strings.Replace(redirect, "%2f", "/", -1)
	http.Redirect(res, req, "/"+redirect, http.StatusSeeOther)
}

// Logins using OAuth
func LoginFromOauth(res http.ResponseWriter, req *http.Request, email string) error {
	ctx := NewContext(res,req)
	l := LoginOauthAccount{}
	err := retrievable.GetEntity(ctx, email, &l)
	if err != nil { return ERROR_NoUser }
	sessID, err := AUTH_CreateSessionID(ctx, l.UserID)
	if err != nil { return err }
	err = COOKIE_Make(res, "session", strconv.FormatInt(sessID, 10))
	if err != nil { return err }
	return nil
}

// Registers using OAuth
func RegisterFromOauth(res http.ResponseWriter, req *http.Request, email, first, last string) error {
	ctx := NewContext(res,req)
	checkLogin := LoginOauthAccount{}

	// Check that user does not exist
	if checkErr := retrievable.GetEntity(ctx, email, &checkLogin); checkErr == nil { return checkErr }
	u := USER_User{
		Email: email,
		First: first,
		Last:  last,
	}
	ukey, putUserErr := retrievable.PlaceEntity(ctx, int64(0), &u)
	if putUserErr != nil { return putUserErr }
	uLogin := LoginOauthAccount{}
	uLogin.UserID = ukey.IntID()
	lkey, putErr := retrievable.PlaceEntity(ctx, email, &uLogin)
	if putErr != nil { return putErr }
	sessID, err := AUTH_CreateSessionID(ctx, lkey.IntID())
	if err != nil { return err }
	err = COOKIE_Make(res, "session", strconv.FormatInt(sessID, 10))
	if err != nil { return err }
	return nil
}
