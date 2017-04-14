package PM

import(
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/USERS"
)

func TestSendMessage(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in SendMessage")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{Email:"test1@test1",})
	retrievable.PlaceEntity(ctx,int64(2),&USERS.User{Email:"test2@test2",})
	retrievable.PlaceEntity(ctx,"test1@test1",&AUTH.EmailReference{int64(1)})
	retrievable.PlaceEntity(ctx,"test2@test2",&AUTH.EmailReference{int64(2)})	
	// Make Context
	UserCtx := CONTEXT.Context{}
	UserCtx.User = &USERS.User{Email:"test1@test1",IntID:retrievable.IntID(1),}
	UserCtx.Context  = ctx
	
	// User1 Sends Message to User2
	SendMessage(UserCtx,"test2@test2","TestTitle","Test")
	
	// Assert that header exists for both.
	Header := PrivateMessageHeader{}
	if retrievable.GetEntity(ctx,int64(1),&Header) != nil {
		fmt.Println("FAIL SendMessage 1")
		t.Fail()
	}
	if retrievable.GetEntity(ctx,int64(2),&Header) != nil {
		fmt.Println("FAIL SendMessage 2")
		t.Fail()
	}
	// Assert that the sent message actually exists.
	if len(Header.Messages) == 1 {
		Message := PrivateMessage{}
		retrievable.GetEntity(ctx,Header.Messages[0],&Message)
		if Message.Content != "Test" || Message.Sender != "test1@test1" {
			fmt.Println("FAIL SendMessage 3",Message)
			t.Fail()
		}
	}
}

func TestRetrieveMessages(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RetrieveMessages")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{Email:"test1@test1",})
	retrievable.PlaceEntity(ctx,int64(2),&USERS.User{Email:"test2@test2",})
	retrievable.PlaceEntity(ctx,"test1@test1",&AUTH.EmailReference{int64(1)})
	retrievable.PlaceEntity(ctx,"test2@test2",&AUTH.EmailReference{int64(2)})	
	// Make Context
	UserCtx := CONTEXT.Context{}
	UserCtx.User = &USERS.User{Email:"test1@test1",IntID:retrievable.IntID(1),}
	UserCtx.Context  = ctx
	// Test With Empty Header
	t1 := RetrieveMessages(UserCtx,2,0)
	if len(t1) != 0 {
		fmt.Println("FAIL RetrieveMessages 1")
		t.Fail()
	}
	// Test With 1 Entry
	SendMessage(UserCtx,"test2@test2","1","Test")
	t2 := RetrieveMessages(UserCtx,2,0)
	if t2[0].Title != "1" {
		fmt.Println("FAIL RetrieveMessages 2")	
		t.Fail()
	}
	// Test Out of Bounds
	t3 := RetrieveMessages(UserCtx,2,1)
	if len(t3) != 0 {
		fmt.Println("FAIL RetrieveMessages 3")
		t.Fail()
	}
	// Test With 3 Entries
	SendMessage(UserCtx,"test2@test2","2","Test")
	SendMessage(UserCtx,"test2@test2","3","Test")
	t4 := RetrieveMessages(UserCtx,2,0)
	if t4[0].Title != "3" || t4[1].Title != "2" {
		fmt.Println("FAIL RetrieveMessages 4",t4)
		t.Fail()
	}
	// Test End of 6 Entries
	t5 := RetrieveMessages(UserCtx,2,1)
	if len(t5) != 1 || t5[0].Title != "1" {
		fmt.Println("FAIL RetrieveMessages 5",t5)
		t.Fail()
	}
}

func TestGetPageNumbers(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetPageNumbers")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{})
	
	// Make Context
	UserCtx := CONTEXT.Context{}
	UserCtx.User = &USERS.User{Email:"test1@test1",IntID:retrievable.IntID(1),}
	UserCtx.Context  = ctx
	
	// 0 element case
	retrievable.PlaceEntity(ctx,int64(1),&PrivateMessageHeader{Messages:[]int64{},})
	if len(GetPageNumbers(UserCtx,5,0)) != 0 {
		fmt.Println("FAIL GetPageNumbers 1")
		t.Fail()
	}
	// 1 element case
	retrievable.PlaceEntity(ctx,int64(1),&PrivateMessageHeader{Messages:[]int64{0},})
	if len(GetPageNumbers(UserCtx,5,0)) != 1 {
		fmt.Println("FAIL GetPageNumbers 2")
		t.Fail()	
	}
	// 7 element cases
	retrievable.PlaceEntity(ctx,int64(1),&PrivateMessageHeader{Messages:[]int64{0,0},})
	// Beginning Case
	t1 := GetPageNumbers(UserCtx,5,0)
	if t1[0] != 0 || t1[1] != 1 || t1[2] != 2 || t1[3] != 3 || t1[4] != 4 {
		fmt.Println("FAIL GetPageNumbers 3")
		t.Fail()	
	}
	// Middle Case
	t2 := GetPageNumbers(UserCtx,5,3)
	if t2[0] != 1 || t2[1] != 2 || t2[2] != 3 || t2[3] != 4 || t2[4] != 5 {
		fmt.Println("FAIL GetPageNumbers 4")
		t.Fail()	
	}
	// End Case
	t3 := GetPageNumbers(UserCtx,5,6)
	if t3[0] != 2 || t3[1] != 3 || t3[2] != 4 || t3[3] != 5 || t3[4] != 6 {
		fmt.Println("FAIL GetPageNumbers 5")
		t.Fail()	
	}
}