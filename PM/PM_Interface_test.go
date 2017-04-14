package PM

import(
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
)

func TestSendMessage(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RenameFolder")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{Email:"test1@test1",})
	retrievable.PlaceEntity(ctx,int64(2),&USERS.User{Email:"test2@test2",})
	
	// Make Context
	UserCtx := CONTEXT.Context{}
	UserCtx.User = &USERS.User{Email:"test1@test1",IntID:retrievable.IntID(1),}
	UserCtx.Ctx  = ctx
	
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
		retrievable.GetEntity(ctx,Header.Messages[0],&PrivateMessage{})
		if Message.Content != "Test" || Message.Sender != "test1@test1" {
			fmt.Println("FAIL SendMessage 3")
			t.Fail()
		}
	}
}
/*
func TestRetrieveMessages(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RenameFolder")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,KEY,&OBJECT)
}

func TestGetPageNumbers(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RenameFolder")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,KEY,&OBJECT)
}
*/