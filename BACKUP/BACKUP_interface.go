package BACKUP
import (
	"golang.org/x/net/context"
	"strconv"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/NOTES"
)
func UpdateBackup(ctx context.Context,NoteID,UserID int64,backup *Backup) error {	
	_ , err := retrievable.PlaceEntity(ctx,strconv.FormatInt(NoteID,10)+"-"+strconv.FormatInt(UserID,10),backup)
	return err
}
func RetrieveBackup(ctx context.Context,NoteID,UserID int64) NOTES.Content { 
	b := Backup{}
	if retrievable.GetEntity(ctx,strconv.FormatInt(NoteID,10)+"-"+strconv.FormatInt(UserID,10),&b) != nil { return NOTES.Content{} }
	return NOTES.Content(b)
}
func BackupExists(ctx context.Context,NoteID,UserID int64) bool {
	return (retrievable.GetEntity(ctx,strconv.FormatInt(NoteID,10)+"-"+strconv.FormatInt(UserID,10),&Backup{}) == nil)
}