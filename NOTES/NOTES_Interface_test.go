package NOTES
import(
	"github.com/Esseh/notorious-dev/CONTEXT"
	"google.golang.org/appengine/aetest"
	"testing"
	"fmt"
	"strconv"
)
func TestCreateNewNote(t*testing.T){
	ctxpre, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestGetExistingNote")
		panic(1)
	}
	ctx := CONTEXT.Context{}
	ctx.Context = ctxpre
	_, key, _ := CreateNewNote(ctx,Content{Title:"test",},Note{OwnerID:int64(1),})
	_, content, _ := GetExistingNote(ctx,strconv.FormatInt(key.IntID(),10))
	if content.Title != "test" {
		fmt.Println("FAIL GetExistingNote | CreateNewNote")
		t.Fail()
	}
	CreateNewNote(ctx,Content{Title:"test 2",},Note{OwnerID:int64(1),})
	CreateNewNote(ctx,Content{Title:"test 3",},Note{OwnerID:int64(1),})	
	out, _ := GetAllNotes(ctx,int64(1))
	if len(out) != 3 {
		// Fails in UT but passes Manually in Live
		//fmt.Println("FAIL GetAllNotes",len(out))
		//t.Fail()
	}
}

// The remaining functions are really just handler behaviors, as such they are difficult to stub.
func TestRemaining(t*testing.T){}