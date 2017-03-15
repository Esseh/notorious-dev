package BACKUP
import(
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"github.com/Esseh/notorious-dev/NOTES"
)
var BackupTable = "Backups"

type Backup NOTES.Content
func (b *Backup) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, BackupTable, key.(string), int64(0), nil)
}