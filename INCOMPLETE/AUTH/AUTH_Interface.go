package AUTH

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"bytes"
	"crypto/aes"
	"encoding/base64"
	
	"github.com/Esseh/retrievable"
	"github.com/mssola/user_agent"
	"github.com/pariz/gountries"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)


// Retrieves an ID for AUTH_User from login information.
func AUTH_GetUserIDFromLogin(ctx Context, email, password string) (int64, error) {
	urID := AUTH_LoginLocalAccount{}
	if getErr := retrievable.GetEntity(ctx, email, &urID); getErr != nil { return -1, getErr }
	if compareErr := bcrypt.CompareHashAndPassword(urID.Password, []byte(password)); compareErr != nil {
		return -1, compareErr
	}
	return urID.UserID, nil
}

// Utilizing an AUTH_User and username/password information it creates a database entry for their AUTH_LoginLocalAccount.
func AUTH_CreateUserFromLogin(ctx Context, email, password string, u *USER_User) (*USER_User, error) {
	checkLogin := AUTH_LoginLocalAccount{}
	// Check that user does not exist
	if checkErr := retrievable.GetEntity(ctx, email, &checkLogin); checkErr == nil {
		return u, ERROR_UsernameExists
	} else if checkErr != datastore.ErrNoSuchEntity && checkErr != nil {
		return u, checkErr
	}

	ukey, putUserErr := retrievable.PlaceEntity(ctx, retrievable.IntID(0), u)
	if putUserErr != nil { return u, putUserErr }
	if u.IntID == 0 { return u, errors.New("HEY, DATASTORE IS STUPID") }

	cryptPass, cryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if cryptErr != nil { return u, cryptErr }

	uLogin := AUTH_LoginLocalAccount{
		Password: cryptPass,
		UserID:   ukey.IntID(),
	}
	_, putErr := retrievable.PlaceEntity(ctx, email, &uLogin)
	return u, putErr
}

// Initializes a new AUTH_Session and returns the ID of that AUTH_Session.
func AUTH_CreateSessionID(ctx Context, userID int64) (sessionID int64, _ error) {
	agent := user_agent.New(ctx.req.Header.Get("user-agent"))
	browse, vers := agent.Browser()
	ip, _, err := net.SplitHostPort(ctx.req.RemoteAddr)
	if err != nil { ip = ctx.req.RemoteAddr }
	country := ctx.req.Header.Get("X-AppEngine-Country")
	region := ctx.req.Header.Get("X-AppEngine-Region")
	city := ctx.req.Header.Get("X-AppEngine-City")
	location, err := AUTH_GetLocationName(country, strings.ToUpper(region))
	if err != nil {
		location = "Unknown"
	} else {
		location = strings.Title(city) + ", " + location
	}
	newSession := AUTH_Session{
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
func AUTH_LoginToWebsite(ctx Context,username,password string) (string, error) {
	userID, err := AUTH_GetUserIDFromLogin(ctx, strings.ToLower(username), password)
	if err != nil { return "Login Information Is Incorrect", err }
	sessionID, err := AUTH_CreateSessionID(ctx, userID)
	if err != nil { return "Login error, try again later.", err }
	err = COOKIE_Make(ctx.res, "session", strconv.FormatInt(sessionID, 10))
	return "Login error, try again later.",err
}

// Makes the currently active user log out.
func AUTH_LogoutFromWebsite(ctx Context)(string, error){
	sessionIDStr, err := COOKIE_GetValue(ctx.req, "session")
	if err != nil { return "Must be logged in", err }
	sessionVal, err := strconv.ParseInt(sessionIDStr, 10, 0)	
	if err != nil { return "Bad cookie value", err }
	err = retrievable.DeleteEntity(ctx, (&AUTH_Session{}).Key(ctx, sessionVal))
	if err == nil { COOKIE_Delete(ctx.res, "session") }
	return "No such session found!", err
}

// Registers a user with the following information...
//	username
//	password
//	confirmPassword
//	firstName
//	lastName
func AUTH_RegisterNewUser(ctx Context, username, password, confirmPassword, firstName, lastName string)(string,error){
	newUser := &USER_User{ // Make the New User
		Email:    strings.ToLower(username),
		First:    firstName,
		Last:     lastName,
	}		
	if !AUTH_ValidLogin(username,password) { return "Invalid Login Information", errors.New("Bad Login") }
	if password != confirmPassword { return "Passwords Do Not Match", errors.New("Password Mismatch") }
	_, err := AUTH_CreateUserFromLogin(ctx, newUser.Email, password, newUser)
	return "Username Taken", err
}