package FORUM

import(
	"net/url"
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"net/http/httptest"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
)

func TestRegisterAdminAPI(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RegisterAdminAPI")
		panic(1)
	}
	
	// Make ctx1
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("Password","password")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Make ctx2
	req2 := httptest.NewRequest("GET", "/", nil)
	values2 := url.Values{}; 
	values2.Add("Password","wordpass")
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx; 
	ctx2.Req = req2
	ctx2.User = &USERS.User{IntID: retrievable.IntID(1),}
	
	// Success Case
	if RegisterAdminAPI(ctx1) != `{"success":true}` {
		t.Fail()
		fmt.Println("Fail FORUM API 1")
	}
	
	// Fail Case
	if RegisterAdminAPI(ctx2) != `{"success":false}` {
		t.Fail()
		fmt.Println("Fail FORUM API 2")	
	}	
	
	// Assert Success
	ah := AdminHeader{}
	retrievable.GetEntity(ctx,nil,&ah)
	if len(ah.UserIDS) != 1 || ah.UserIDS[0] != 1 {
		t.Fail()
		fmt.Println("Fail FORUM API 3")
	}
}