package USERS

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"
	"strconv"
	"mime/multipart"
	"errors"
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/CLOUD"
)


// Uploads an avatar for a user into the cloud.
func UploadAvatar(ctx context.Context, userID int64, header *multipart.FileHeader, avatarReader io.ReadSeeker) error {
	m, _, err := image.Decode(avatarReader)
	if err != nil { return err }
	imageBounds := m.Bounds()
	if imageBounds.Dy() > CORE.MaxAvatarSize || imageBounds.Dx() > CORE.MaxAvatarSize {
		return errors.New("Filesize is Too Large")
	}
	avatarReader.Seek(0, 0)
	filename := CORE.GetAvatarPath(userID)
	return CLOUD.AddFile(ctx, filename, header.Header["Content-Type"][0], avatarReader)
}

// Retrieves an AUTH_User from the currently logged in user.
func GetUserFromSession(ctx context.Context,req *http.Request) (*User, error) {
	// Get session ID from cookie
	sessionIDString, err := COOKIE.GetValue(req, "session")
	if err != nil { return &User{}, err) }
	sessionID, _ := strconv.ParseInt(sessionIDString, 10, 64) // Change cookie val into key	

	// get session data
	session := Session{}
	err := retrievable.GetEntity(ctx,sessionID,&session)
	if err != nil { return &User{}, err }
	
	// get user id from session data
	userID := session.UserID
	
	if userID != 0 {	
		// Get User
		user := User{}
		retrievable.GetEntity(ctx,userID,&user)

		// Update Session Information
		session.LastUsed = time.Now()
		retrievable.PlaceEntity(ctx,sessionID,&session)

		// Return Result
		return &user, nil	
	} else {
		return &User{}, errors.New("0 Key is Bad Key")
	}
}