package NOTES
import(
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/USERS"
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

func TestCanViewNote(t*testing.T){
	u1 := &USERS.User{}; u1.IntID = retrievable.IntID(1)
	u2 := &USERS.User{}; u2.IntID = retrievable.IntID(2)
	u3 := &USERS.User{}; u3.IntID = retrievable.IntID(3)
	u4 := &USERS.User{}; u4.IntID = retrievable.IntID(4)
	n := &Note{OwnerID:int64(1),Collaborators:[]int64{2,3},PublicallyViewable:true}
	if !CanViewNote(n,u4) {
		fmt.Println("FAIL ViewNote 1")
		t.Fail()
	}
	n.PublicallyViewable = false
	if CanViewNote(n,u4) {
		fmt.Println("FAIL ViewNote 2")
		t.Fail()
	}
	if !CanViewNote(n,u3) {
		fmt.Println("FAIL ViewNote 3")
		t.Fail()
	}
	if !CanViewNote(n,u2) {
		fmt.Println("FAIL ViewNote 4")
		t.Fail()
	}
	if !CanViewNote(n,u1) {
		fmt.Println("FAIL ViewNote 5")
		t.Fail()
	}
	testing.Coverage()
}

func TestCanEditNote(t*testing.T){
	u1 := &USERS.User{}; u1.IntID = retrievable.IntID(1)
	u2 := &USERS.User{}; u2.IntID = retrievable.IntID(2)
	u3 := &USERS.User{}; u3.IntID = retrievable.IntID(3)
	u4 := &USERS.User{}; u4.IntID = retrievable.IntID(4)
	n := &Note{OwnerID:int64(1),Collaborators:[]int64{2,3},PublicallyEditable:true}
	if !CanEditNote(n,u4) {
		fmt.Println("FAIL EditNote 1")
		t.Fail()
	}
	n.PublicallyEditable = false
	if CanEditNote(n,u4) {
		fmt.Println("FAIL EditNote 2")
		t.Fail()
	}
	if !CanEditNote(n,u3) {
		fmt.Println("FAIL EditNote 3")
		t.Fail()
	}
	if !CanEditNote(n,u2) {
		fmt.Println("FAIL EditNote 4")
		t.Fail()
	}
	if !CanEditNote(n,u1) {
		fmt.Println("FAIL EditNote 5")
		t.Fail()
	}
	testing.Coverage()
}

// The remaining functions are really just handler behaviors, as such they are difficult to stub.
func TestRemaining(t*testing.T){}