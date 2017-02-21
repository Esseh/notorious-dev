package USERS
import(
	"testing"
	"fmt"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/notorious-dev/COOKIE"
	"net/http/httptest"
	"github.com/Esseh/retrievable"
)

// func TestUploadAvatar(t *testing.T)
// I cannot find a way to properly stub MIME Encoding with Go.
// Here is an example of proper usage...
/*
	rdr, hdr, _ := req.FormFile("avatar")
	defer rdr.Close()
	USERS.UploadAvatar(..., ..., hdr, rdr)
}*/

func TestGetUserFromSession(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestUploadAvatar")
		panic(1)
	}
	res := httptest.NewRecorder()
	COOKIE.Make(res,"session","1")
	req := httptest.NewRequest("GET", "/", nil)
	retrievable.PlaceEntity(ctx,int64(1),&User{First:"test",})
	retrievable.PlaceEntity(ctx,int64(1),&Session{UserID:int64(1),})
	if _ , err := GetUserFromSession(ctx,req); err == nil {
		fmt.Println("FAIL GetUserFromSession 1")
		t.Fail()
	}
	req.AddCookie(res.Result().Cookies()[0])
	if u, _ := GetUserFromSession(ctx,req); u.First != "test" {
		fmt.Println("FAIL GetUserFromSession 2")
		t.Fail()	
	}
	testing.Coverage()
}