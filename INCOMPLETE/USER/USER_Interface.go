package USER

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"strconv"
	"golang.org/x/net/context"
)


// Uploads an avatar for a user into the cloud.
func USER_UploadAvatar(ctx context.Context, userID int64, header *multipart.FileHeader, avatarReader io.ReadSeeker) error {
	m, _, err := image.Decode(avatarReader)
	if err != nil { return err }
	imageBounds := m.Bounds()
	if imageBounds.Dy() > maxAvatarSize || imageBounds.Dx() > maxAvatarSize {
		return ERROR_TooLarge
	}
	if _, err = avatarReader.Seek(0, 0); err != nil {
		return err
	}
	filename := getAvatarPath(userID)
	return CLOUD_AddFile(ctx, filename, header.Header["Content-Type"][0], avatarReader)
}