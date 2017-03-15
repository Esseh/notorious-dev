package BACKUP
import (
	"golang.org/x/net/context"
	"github.com/Esseh/notorious-dev/NOTES"
)
func UpdateBackup(ctx context.Context,NoteID,UserID int64,backup Backup) error { return nil}
func RetrieveBackup(ctx context.Context,NoteID,UserID int64) NOTES.Content { return NOTES.Content{}}
func BackupExists(ctx context.Context,NoteID,UserID int64) bool { return false }