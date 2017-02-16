package AUTH

import (
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const (
	LoginTable         = "Login"
)

type (
	// Contains local authentication information about the user.
	LoginLocalAccount struct {
		// Key for USER_User
		UserID   int64
		// An encrypted password.
		Password []byte
	}

)

// Retrieves the Local Login Account, overloaded to handle OAUTH.
func (l *LoginLocalAccount) Get(ctx context.Context, key interface{}) error {
	getErr := retrievable.GetEntity(ctx, key, l) 
	if getErr != nil {                           
		oauth := LoginOauthAccount{}
		ogetErr := retrievable.GetEntity(ctx, key, &oauth)
		if ogetErr != nil { return ogetErr }
		l.UserID = oauth.UserID
	}
	return nil
}

//Updates the Local Login Account, overloaded to handle OAUTH.
func (l *LoginLocalAccount) Place(ctx context.Context, key interface{}) (*datastore.Key, error) {
	if string(l.Password) == "" { // OAuth Case
		oauth := LoginOauthAccount{}
		oauth.UserID = oauth.UserID
		return retrievable.PlaceEntity(ctx, key, &oauth)
	} else { // LoginLocal Case
		return retrievable.PlaceEntity(ctx, key, l)
	}
}

func (l *LoginLocalAccount) Key(ctx context.Context, key interface{}) *datastore.Key {
	e, _ := CORE.Encrypt([]byte(key.(string)), encryptKey)
	return datastore.NewKey(ctx, LoginTable, e, 0, nil)
}


const ( 
	OauthTable = "Oauth" 
)
// Login Information for Oauth. Functionally equivalent to LoginLocalAccount
type LoginOauthAccount struct {
	UserID int64
}
// String Keys
func (l *LoginOauthAccount) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, OauthTable, key.(string), 0, nil)
}