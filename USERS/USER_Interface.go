package USERS

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"errors"
	"golang.org/x/net/context"
	"github.com/Esseh/notorious-dev/CORE"
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
	if _, err = avatarReader.Seek(0, 0); err != nil {
		return err
	}
	filename := CORE.GetAvatarPath(userID)
	return CLOUD.AddFile(ctx, filename, header.Header["Content-Type"][0], avatarReader)
}