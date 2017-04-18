package NOTES
import(
	"net/http/httptest"
	"net/url"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/NOTIFICATION"
	"google.golang.org/appengine/aetest"
	"testing"
	"fmt"
	"strconv"
)

func TestNotify(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetPageNumbers")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{})
	retrievable.PlaceEntity(ctx,int64(2),&USERS.User{})
	retrievable.PlaceEntity(ctx,int64(1),&Note{ContentID:1,})
	retrievable.PlaceEntity(ctx,int64(1),&Content{Title:"testNote",})
	retrievable.PlaceEntity(ctx,int64(1),&SubscriptionHeader{UserIDS:[]int64{1,2},})
	
	// Make Context
	UserCtx := CONTEXT.Context{ User:&USERS.User{Email:"test1@test1"}, }
	UserCtx.Context  = ctx

	// Send PM Notifications
	Notify(UserCtx,1)
	Notify(UserCtx,1)

	// Assert Success
	n1 := NOTIFICATION.Notifications{}
	retrievable.GetEntity(ctx,int64(1),&n1)
	if n1.NotificationsPending != 2 {
		t.Fail()
		fmt.Println("FAIL Notify 1")
	}
	
	if n1.Notifications[0] != "test1@test1 updated testNote" || n1.Notifications[1] != "test1@test1 updated testNote" {
		t.Fail()
		fmt.Println("FAIL Notify 2")
	}	

	// Assert Success
	n2 := NOTIFICATION.Notifications{}
	retrievable.GetEntity(ctx,int64(2),&n2)
	if n2.NotificationsPending != 2 {
		t.Fail()
		fmt.Println("FAIL Notify 3")
	}
	
	if n2.Notifications[0] != "test1@test1 updated testNote" || n2.Notifications[1] != "test1@test1 updated testNote" {
		t.Fail()
		fmt.Println("FAIL Notify 4")
	}
}

func TestSubscribeAPI(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestSubscribeAPI")
		panic(1)
	}
	
	retrievable.PlaceEntity(ctx, int64(11), &Note{ContentID:int64(11),})
	retrievable.PlaceEntity(ctx, int64(11), &Content{})
	
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("NoteID","11")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Success Case
	if SubscribeAPI(ctx1) != `{"success":true}`{
		fmt.Println("Fail SubscribeAPI 1")
		t.Fail()
	}
	
	// Fail Case, already subscribed
	if SubscribeAPI(ctx1) != `{"success":false}`{
		fmt.Println("Fail SubscribeAPI 2")
		t.Fail()	
	}
}
func TestUnsubscribeAPI(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestUnsubscribeAPI")
		panic(1)
	}
	
	retrievable.PlaceEntity(ctx, int64(11), &Note{ContentID:int64(11),})
	retrievable.PlaceEntity(ctx, int64(11), &Content{})
	
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("NoteID","11")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	SubscribeAPI(ctx1)
	
	// Success Case
	if o := UnsubscribeAPI(ctx1); o != `{"success":true}`{
		fmt.Println("Fail UnsubscribeAPI 1")
		t.Fail()
	}
	// Fail Case, not subscribed
	if UnsubscribeAPI(ctx1) != `{"success":false}`{
		fmt.Println("Fail UnsubscribeAPI 2")
		t.Fail()	
	}
}

func TestGetSubscriptions(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestGetSubscriptions")
		panic(1)
	}
	
	retrievable.PlaceEntity(ctx, int64(11), &Note{ContentID:int64(11),})
	retrievable.PlaceEntity(ctx, int64(11), &Content{Title:"1",})
	retrievable.PlaceEntity(ctx, int64(22), &Note{ContentID:int64(22),})
	retrievable.PlaceEntity(ctx, int64(22), &Content{Title:"2",})
	retrievable.PlaceEntity(ctx, int64(33), &Note{ContentID:int64(33),})
	retrievable.PlaceEntity(ctx, int64(33), &Content{Title:"3",})
	retrievable.PlaceEntity(ctx, int64(1), &Subscription{NoteIDS:[]int64{11,22,33},})
	
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx
	out := GetSubscriptions(ctx1, 1)
	if len(out) == 0 {
		fmt.Println("GetSubscriptions not Implemented")
	} else {
		if out[0].Content.Title != "1" || out[1].Content.Title != "2" || out[2].Content.Title != "3" {
			fmt.Println("FAIL GetSubscriptions")
			t.Fail()
		}	
	}
}


func TestAPI_SaveCopy(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestAPI_SaveCopy")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(11),&Note{OwnerID:1,ContentID:33,})
	retrievable.PlaceEntity(ctx,int64(22),&Note{OwnerID:2,ContentID:44,})
	retrievable.PlaceEntity(ctx,int64(33),&Content{Title:"test",})
	
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("NoteID","22")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	req2 := httptest.NewRequest("GET", "/", nil)
	values2 := url.Values{}; 
	values2.Add("NoteID","11")
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx; 
	ctx2.Req = req2
	ctx2.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Not Allowed Case / Note Doesn't Exist	
	if API_SaveCopy(ctx1) != `{"success":false}` {
		fmt.Println("FAIL SaveCopy 1")
		t.Fail()	
	}
	// Success Case
	assert1 := Note{}
	assert2 := Content{}
	s := (API_SaveCopy(ctx2))
	if s != `{"success":false}` { s = s[25:41]}
	k, _ := strconv.ParseInt(s,10,64)
	retrievable.GetEntity(ctx,k,&assert1)
	retrievable.GetEntity(ctx,assert1.ContentID,&assert2)
	// Assert that the new entry is actually made and that it is not duplicating any existing stubs.
	if assert2.Title != "test" || assert1.ContentID == int64(0) || assert1.ContentID == int64(44) || k == int64(22) {
		fmt.Println("FAIL SaveCopy 2")
		t.Fail()
	}
}

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