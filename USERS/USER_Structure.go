// Contains the object structure and methods relating to the USER.
package USERS

import (
	"encoding/json"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	UsersTable          = "Users"
	RecentlyViewedTable = "RecentlyViewedCourses"
)

type (
	// Represents an individual user.
	User struct {
		// First and Last Name
		First, Last       string
		Email             string
		// Whether they have an active avatar or not.
		Avatar            bool `datastore:",noindex"`
		// Biography.
		Bio               string
		// ID referred to itself.
		retrievable.IntID `datastore:"-" json:"-"`
	}
	// An encrypted user.
	EncryptedUser struct {
		First, Last string
		Email       string
		Avatar      bool `datastore:",noindex"`
		Bio         string
	}
)

func (u *User) Key(ctx context.Context, key interface{}) *datastore.Key {
	if v, ok := key.(retrievable.IntID); ok {
		return datastore.NewKey(ctx, UsersTable, "", int64(v), nil)
	}
	return datastore.NewKey(ctx, UsersTable, "", key.(int64), nil)
}

// Converts user to an encrypted user.
func (u *User) toEncrypt() (*EncryptedUser, error) {
	e := &EncryptedUser{
		First:     u.First,
		Last:      u.Last,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
	}
	email, err := CORE.Encrypt([]byte(u.Email), CORE.EncryptKey)
	if err != nil { return nil, err }
	e.Email = email
	return e, nil
}

// Converts encrypted user to normal user.
func (u *User) fromEncrypt(e *EncryptedUser) error {
	email, err := CORE.Decrypt(e.Email, CORE.EncryptKey)
	if err != nil { return err }
	u.First = e.First
	u.Last = e.Last
	u.Email = string(email)
	u.Avatar = e.Avatar
	u.Bio = e.Bio
	return nil
}

// User -> JSON
func (u *User) Serialize() []byte {
	data, _ := u.toEncrypt()
	ret, _ := json.Marshal(&data)
	return ret
}
// JSON -> User
func (u *User) Unserialize(data []byte) error {
	e := &EncryptedUser{}
	err := json.Unmarshal(data, e)
	if err != nil { return err }
	return u.fromEncrypt(e)
}

