package AUTH

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/mssola/user_agent"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine/datastore"
)


// Retrieves an ID for AUTH_User from login information.
func GetUserIDFromLogin(ctx CONTEXT.Context, email, password string) (int64, error) {
	urID := LoginLocalAccount{}
	if getErr := retrievable.GetEntity(ctx, email, &urID); getErr != nil { return -1, getErr }
	if compareErr := bcrypt.CompareHashAndPassword(urID.Password, []byte(password)); compareErr != nil {
		return -1, compareErr
	}
	return urID.UserID, nil
}

// Utilizing an AUTH_User and username/password information it creates a database entry for their AUTH_LoginLocalAccount.
func CreateUserFromLogin(ctx CONTEXT.Context, email, password string, u *USERS.User) (*USERS.User, error) {
	checkLogin := LoginLocalAccount{}
	// Check that user does not exist
	if checkErr := retrievable.GetEntity(ctx, email, &checkLogin); checkErr == nil {
		return u, errors.New("Username Already Exists")
	} else if checkErr != datastore.ErrNoSuchEntity && checkErr != nil {
		return u, checkErr
	}

	ukey, putUserErr := retrievable.PlaceEntity(ctx, retrievable.IntID(0), u)
	if putUserErr != nil { return u, putUserErr }
	if u.IntID == 0 { return u, errors.New("HEY, DATASTORE IS STUPID") }

	cryptPass, cryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if cryptErr != nil { return u, cryptErr }

	uLogin := LoginLocalAccount{
		Password: cryptPass,
		UserID:   ukey.IntID(),
	}
	_, putErr := retrievable.PlaceEntity(ctx, email, &uLogin)
	return u, putErr
}

// Initializes a new AUTH_Session and returns the ID of that AUTH_Session.
func CreateSessionID(ctx CONTEXT.Context, userID int64) (sessionID int64, _ error) {
	agent := user_agent.New(ctx.req.Header.Get("user-agent"))
	browse, vers := agent.Browser()
	ip, _, err := net.SplitHostPort(ctx.req.RemoteAddr)
	if err != nil { ip = ctx.req.RemoteAddr }
	country := ctx.req.Header.Get("X-AppEngine-Country")
	region := ctx.req.Header.Get("X-AppEngine-Region")
	city := ctx.req.Header.Get("X-AppEngine-City")
	location, err := CORE.GetLocationName(country, strings.ToUpper(region))
	if err != nil {
		location = "Unknown"
	} else {
		location = strings.Title(city) + ", " + location
	}
	newSession := USERS.Session{
		UserID:      userID,
		BrowserUsed: browse + " " + vers,
		IP:          ip,
		LocationUsed: location,
		LastUsed:     time.Now(),
	}
	rk, err := retrievable.PlaceEntity(ctx, int64(0), &newSession)
	if err != nil { return int64(-1), err }
	return rk.IntID(), err
}


// Makes the currently active user log in with username and password information.
func LoginToWebsite(ctx CONTEXT.Context,username,password string) (string, error) {
	userID, err := GetUserIDFromLogin(ctx, strings.ToLower(username), password)
	if err != nil { return "Login Information Is Incorrect", err }
	sessionID, err := CreateSessionID(ctx, userID)
	if err != nil { return "Login error, try again later.", err }
	err = COOKIE.Make(ctx.res, "session", strconv.FormatInt(sessionID, 10))
	return "Login error, try again later.",err
}

// Makes the currently active user log out.
func LogoutFromWebsite(ctx CONTEXT.Context)(string, error){
	sessionIDStr, err := COOKIE.GetValue(ctx.req, "session")
	if err != nil { return "Must be logged in", err }
	sessionVal, err := strconv.ParseInt(sessionIDStr, 10, 0)	
	if err != nil { return "Bad cookie value", err }
	err = retrievable.DeleteEntity(ctx, (&USERS.Session{}).Key(ctx, sessionVal))
	if err == nil { COOKIE.Delete(ctx.res, "session") }
	return "No such session found!", err
}

// Registers a user with the following information...
//	username
//	password
//	confirmPassword
//	firstName
//	lastName
func RegisterNewUser(ctx CONTEXT.Context, username, password, confirmPassword, firstName, lastName string)(string,error){
	newUser := &USERS.User{ // Make the New User
		Email:    strings.ToLower(username),
		First:    firstName,
		Last:     lastName,
	}		
	if !CORE.ValidLogin(username,password) { return "Invalid Login Information", errors.New("Bad Login") }
	if password != confirmPassword { return "Passwords Do Not Match", errors.New("Password Mismatch") }
	_, err := CreateUserFromLogin(ctx, newUser.Email, password, newUser)
	return "Username Taken", err
}

// Logs the user in with an OAuth id.
func OAuthLogin(req *http.Request, res http.ResponseWriter, id, first, last, redirect string) {
	err := LoginFromOauth(res, req, id)
	if err == errors.New("There is no existing user.") {
		RegisterFromOauth(res, req, id, first, last)
	}
	redirect = strings.Replace(redirect, "%2f", "/", -1)
	http.Redirect(res, req, "/"+redirect, http.StatusSeeOther)
}

// Logins using OAuth
func LoginFromOauth(res http.ResponseWriter, req *http.Request, email string) error {
	ctx := NewCONTEXT.Context(res,req)
	l := LoginOauthAccount{}
	err := retrievable.GetEntity(ctx, email, &l)
	if err != nil { return errors.New("There is no existing user.") }
	sessID, err := CreateSessionID(ctx, l.UserID)
	if err != nil { return err }
	err = COOKIE.Make(res, "session", strconv.FormatInt(sessID, 10))
	if err != nil { return err }
	return nil
}

// Registers using OAuth
func RegisterFromOauth(res http.ResponseWriter, req *http.Request, email, first, last string) error {
	ctx := NewCONTEXT.Context(res,req)
	checkLogin := LoginOauthAccount{}

	// Check that user does not exist
	if checkErr := retrievable.GetEntity(ctx, email, &checkLogin); checkErr == nil { return checkErr }
	u := USERS.User{
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
	sessID, err := CreateSessionID(ctx, lkey.IntID())
	if err != nil { return err }
	err = COOKIE.Make(res, "session", strconv.FormatInt(sessID, 10))
	if err != nil { return err }
	return nil
}
