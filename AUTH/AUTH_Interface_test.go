package AUTH
import(
	"testing"
	"google.golang.org/appengine/aetest"
	"fmt"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/CONTEXT"
)

func TestGetUserIDFromLogin(t*testing.T){
	ctxpre, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestUploadAvatar")
		panic(1)
	}
	
	ctx := CONTEXT.Context{}
	ctx.Context = ctxpre
	
	_, err = GetUserIDFromLogin(ctx, "test@test", "test")
	if err == nil {
		fmt.Println("FAIL GetUserIDFromLogin 1")
		t.Fail()
	}

	err = CreateUserFromLogin(ctx, "test@test", "test", &USERS.User{})	
	if err != nil {
		fmt.Println("FAIL CreateUserFromLogin 1")
		t.Fail()	
	}
	err = CreateUserFromLogin(ctx, "test@test", "test", &USERS.User{})	
	if err == nil {
		fmt.Println("FAIL CreateUserFromLogin 2")
		t.Fail()		
	}
	
	// POST CreateUserFromLogin
	_, err = GetUserIDFromLogin(ctx, "test@test", "test")
	if err != nil {
		fmt.Println("FAIL GetUserIDFromLogin 2")
		t.Fail()	
	}
	testing.Coverage()
}

// This function does not properly work without data from inflight requests. :(
// Has to be manually tested. Though it shouldn't need to be modified in the future.
// Unfortunantly the remaining functions under AUTH have this function as a dependancy.
// Thankfully manual testing of the remaining functions are trivial as the impacts on the website will be obvious ie: login interactions
func TestCreateSessionID(t*testing.T){}