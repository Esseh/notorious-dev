package BACKUP
import(
	"testing"
	"github.com/Esseh/retrievable"
)
func TestBackup(t*testing.T){
	// Would panic if it didn't implement retrievable.
	(func(r retrievable.Retrievable){})(&Backup{})
}