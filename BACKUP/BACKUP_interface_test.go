package BACKUP
import(
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/NOTES"
)
func TestUpdateBackup(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in UpdateBackupTest")
		panic(1)
	}
	if UpdateBackup(ctx,int64(24),int64(42),Backup{Title:"test"}) != nil {
		fmt.Println("FAIL UpdateBackup 1")
		t.Fail()
	}
	// Make sure not only that it is stored, but that the key schema is followed.
	b := Backup{} 
	retrievable.GetEntity(ctx,"24-42",&b)
	if b.Title != "test" {
		fmt.Println("FAIL UpdateBackup 1")
		t.Fail()
	}
}
func TestRetrieveBackup(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RetrieveBackupTest")
		panic(1)
	}
	// Place an Entry
	UpdateBackup(ctx,int64(24),int64(42),Backup{Title:"test"})
	// Success Case
	if (RetrieveBackup(ctx,int64(24),int64(42))).Title != "test" {
		fmt.Println("FAIL RetrieveBackup 1")
		t.Fail()
	}
	// Fail Case
	if RetrieveBackup(ctx,int64(10),int64(10)) != (NOTES.Content{}) {
		fmt.Println("FAIL RetrieveBackup 2")
		t.Fail()
	}	
}
func TestBackupExists(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in BackupExistsTest")
		panic(1)
	}
	// Fail Case
	if BackupExists(ctx,int64(24),int64(42)){
		fmt.Println("FAIL BackupExists 1")
		t.Fail()
	}
	// Place an Entry
	UpdateBackup(ctx,int64(24),int64(42),Backup{Title:"test"})
	// Success Case
	if !BackupExists(ctx,int64(24),int64(42)){
		fmt.Println("FAIL BackupExists 2")
		t.Fail()
	}
}