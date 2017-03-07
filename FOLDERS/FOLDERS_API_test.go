package FOLDERS

import(
	"net/url"
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"net/http/httptest"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
)

func TestAPI_AddNote(t *testing.T){

}
func TestAPI_RemoveNote(t *testing.T){

}
func TestAPI_InitializeRoot(t *testing.T){

}

func TestAPI_OpenFolder(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in OpenFolder")
		panic(1)
	}
	
	// Stub Database
		// Parent Folder
		retrievable.PlaceEntity(ctx,"1",&Folder{
			IsRoot:true,
			OwnerID:int64(1),
			ChildFolders:[]string{"test_child_folder1","test_child_folder2","test_child_folder3"},
			ChildNotes:[]int64{1111,2222,3333,4444},
			ParentFolder:"",
		})
		// Child Folders
		retrievable.PlaceEntity(ctx,"1/test_child_folder1",&Folder{
			OwnerID:int64(1),
			ParentFolder:"1",
		})
		retrievable.PlaceEntity(ctx,"1/test_child_folder2",&Folder{
			OwnerID:int64(1),
			ParentFolder:"1",
		})
		retrievable.PlaceEntity(ctx,"1/test_child_folder3",&Folder{
			OwnerID:int64(1),
			ParentFolder:"1",
		})
		// Child Notes
		retrievable.PlaceEntity(ctx,int64(1111),&NOTES.Note{
			OwnerID:int64(2),
			PublicallyViewable: true,
			ContentID: int64(1111),
		})
		retrievable.PlaceEntity(ctx,int64(1111),&NOTES.Content{Title:"note1"})
		retrievable.PlaceEntity(ctx,int64(2222),&NOTES.Note{
			OwnerID:int64(2),
			ContentID: int64(2222),
		})
		retrievable.PlaceEntity(ctx,int64(2222),&NOTES.Content{Title:"note2"})
		retrievable.PlaceEntity(ctx,int64(3333),&NOTES.Note{
			OwnerID:int64(2),
			PublicallyViewable: true,
			PublicallyEditable: true,
			ContentID: int64(3333),
		})
		retrievable.PlaceEntity(ctx,int64(3333),&NOTES.Content{Title:"note3"})
		
	// Make ctx1
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; values1.Add("FolderID","1")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	validNotes   := `{"title":"note1","id":1111,"noEdit":true},{"title":"note3","id":3333,"noEdit":false}`
	validFolders := `"test_child_folder1","test_child_folder2","test_child_folder3"`
	if OpenFolder(ctx1) != `{"success":true,"code":-1,folders:[`+validFolders+`],"notes":[`+validNotes+`]}` {
		fmt.Println("FAIL OpenFolder 1",OpenFolder(ctx1))
		t.Fail()
	}
	
	// CHECK TO SEE IF DANGLING NOTE IS DELETED
	lastCheck := Folder{}
	lastErr := retrievable.GetEntity(ctx,"1",&lastCheck)
	if len(lastCheck.ChildNotes) == 4 || lastErr != nil || len(lastCheck.ChildNotes) == 0 {
		fmt.Println("FAIL OpenFolder 2")
		t.Fail()
	}
}

func TestAPI_NewFolder(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in NewFolder")
		panic(1)
	}

	// Create Fake User
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{First:"test",})
	retrievable.PlaceEntity(ctx,int64(1),&USERS.Session{UserID:int64(1),})

	// Test User Owns this Folder
	retrievable.PlaceEntity(ctx,"1",&Folder{IsRoot:true,OwnerID:int64(1),})	
	// Test User Doesn't Own this Folder
	retrievable.PlaceEntity(ctx,"2",&Folder{IsRoot:true,OwnerID:int64(2),})	
	
	// Initialize Each Test Request
	res := httptest.NewRecorder(); COOKIE.Make(res,"session","1")
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.AddCookie(res.Result().Cookies()[0])
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(res.Result().Cookies()[0])
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.AddCookie(res.Result().Cookies()[0])
	req4 := httptest.NewRequest("GET", "/", nil)
	req4.AddCookie(res.Result().Cookies()[0])

	// Make ctx1
	values1 := url.Values{}
	values1.Add("ParentID","2")
	values1.Add("FolderName","test")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx2
	values2 := url.Values{}
	values2.Add("ParentID","1")
	// Should be over 600 characters long
	str := "aaaaaaaaaa"; str = str + str; str = str + str; str = str + str; str = str + str; str = str + str; str = str + str; str = str + str;
	values2.Add("FolderName",str)
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx
	ctx2.Req = req2
	ctx2.User = &USERS.User{IntID: retrievable.IntID(1),}	

	// Make ctx3
	values3 := url.Values{}
	values3.Add("ParentID","3")
	values3.Add("FolderName","test")
	req3.Form = values3
	ctx3 := CONTEXT.Context{}
	ctx3.Context = ctx
	ctx3.Req = req3
	ctx3.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx4
	values4 := url.Values{}
	values4.Add("ParentID","1")
	values4.Add("FolderName","test")
	req4.Form = values4
	ctx4 := CONTEXT.Context{}
	ctx4.Context = ctx
	ctx4.Req = req4
	ctx4.User = &USERS.User{IntID: retrievable.IntID(1),}	
	
	// Not Owner Case	
	if NewFolder(ctx1) != `{"success":false,"code":3}` { 
		fmt.Println("FAIL API_NewFolder 1")
		t.Fail()	
	}
	// Filepath Too Long Case
	if NewFolder(ctx2) != `{"success":false,"code":2}` { 
		fmt.Println("FAIL API_NewFolder 2")
		t.Fail()		
	}
	// Database Failure Case
	if NewFolder(ctx3) != `{"success":false,"code":1}` { 
		fmt.Println("FAIL API_NewFolder 3")
		t.Fail()		
	}
	// Normal Case
	if NewFolder(ctx4) != `{"success":true,"code":-1}` { 
		fmt.Println("FAIL API_NewFolder 4")
		fmt.Println(NewFolder(ctx4))
		t.Fail()		
	}
	// Assert that the created folder actually exists.
	if retrievable.GetEntity(ctx,"1/test",&Folder{}) != nil {
		fmt.Println("FAIL API_NewFolder 5")
		t.Fail()
	}
}


func TestAPI_DeleteFolder(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in NewFolder")
		panic(1)
	}

	// Create Fake User
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{First:"test",})
	retrievable.PlaceEntity(ctx,int64(1),&USERS.Session{UserID:int64(1),})

	// Test User Owns these Folders
	retrievable.PlaceEntity(ctx,"1",&Folder{IsRoot:true,OwnerID:int64(1),ChildFolders:[]string{"test"},})	
	retrievable.PlaceEntity(ctx,"1/test",&Folder{IsRoot:false,OwnerID:int64(1),ParentFolder:"1",ChildFolders:[]string{"test2"},})	
	retrievable.PlaceEntity(ctx,"1/test/test2",&Folder{IsRoot:false,OwnerID:int64(1),ParentFolder:"1/test",})
	// Test User Doesn't Own this Folder
	retrievable.PlaceEntity(ctx,"2",&Folder{IsRoot:true,OwnerID:int64(2),})	
	retrievable.PlaceEntity(ctx,"2/test",&Folder{IsRoot:false,OwnerID:int64(2),ParentFolder:"2",})	
	
	// Initialize Each Test Request
	res := httptest.NewRecorder(); COOKIE.Make(res,"session","1")
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.AddCookie(res.Result().Cookies()[0])
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(res.Result().Cookies()[0])
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.AddCookie(res.Result().Cookies()[0])
	req4 := httptest.NewRequest("GET", "/", nil)
	req4.AddCookie(res.Result().Cookies()[0])

	// Make ctx1
	values1 := url.Values{}
	values1.Add("ParentID","2")
	values1.Add("FolderName","test")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx2
	values2 := url.Values{}
	values2.Add("ParentID","")
	values2.Add("FolderName","1")
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx
	ctx2.Req = req2
	ctx2.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx3
	values3 := url.Values{}
	values3.Add("ParentID","3")
	values3.Add("FolderName","test")
	req3.Form = values3
	ctx3 := CONTEXT.Context{}
	ctx3.Context = ctx
	ctx3.Req = req3
	ctx3.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx4
	values4 := url.Values{}
	values4.Add("ParentID","1")
	values4.Add("FolderName","test")
	req4.Form = values4
	ctx4 := CONTEXT.Context{}
	ctx4.Context = ctx
	ctx4.Req = req4
	ctx4.User = &USERS.User{IntID: retrievable.IntID(1),}	
	
	// Not Owner Case	
	if DeleteFolder(ctx1) != `{"success":false,"code":3}` { 
		fmt.Println("FAIL API_DeleteFolder 1",DeleteFolder(ctx1))
		t.Fail()	
	}
	
	// Is Root Case
	if DeleteFolder(ctx2) != `{"success":false,"code":4}` { 
		fmt.Println("FAIL API_DeleteFolder 2",DeleteFolder(ctx2))
		t.Fail()		
	}
	// Database Failure Case
	if DeleteFolder(ctx3) != `{"success":false,"code":1}` { 
		fmt.Println("FAIL API_DeleteFolder 3")
		t.Fail()		
	}
	// Normal Case
	if DeleteFolder(ctx4) != `{"success":true,"code":-1}` { 
		fmt.Println("FAIL API_DeleteFolder 4",DeleteFolder(ctx4))
		t.Fail()		
	}
	// Assert that the created folder actually exists.
	if retrievable.GetEntity(ctx,"1/test",&Folder{}) == nil {
		fmt.Println("FAIL API_DeleteFolder 5")
		t.Fail()
	}
	// Assert that recursive deletion worked.
	if retrievable.GetEntity(ctx,"1/test/test2",&Folder{}) == nil {
		fmt.Println("FAIL API_DeleteFolder 6")
		t.Fail()
	}
}