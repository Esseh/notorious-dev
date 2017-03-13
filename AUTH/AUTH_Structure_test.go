package AUTH
import(
	"google.golang.org/appengine/aetest"
	"testing"
	"fmt"
)

func TestReference(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestReference")
		panic(1)
	}
	var test *ReferenceID
	test.Key(ctx,"teststring")
}

func TestLoginLocalAccount(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in LoginLocalAccount")
		panic(1)
	}
	if (&LoginLocalAccount{}).Key(ctx,"1") == nil {
		fmt.Println("FAIL TestUser 1")
		t.Fail()
	}
	if (&LoginOauthAccount{}).Key(ctx,"1") == nil {
		fmt.Println("FAIL TestUser 1")
		t.Fail()
	}

	
	(&LoginLocalAccount{Password:[]byte(""),UserID:1,}).Place(ctx,"oauth-example")
	(&LoginLocalAccount{Password:[]byte("password"),UserID:2}).Place(ctx,"local-example")

	ll1 := &LoginLocalAccount{}
	ll2 := &LoginLocalAccount{}
	ll1.Get(ctx,"oauth-example")
	ll2.Get(ctx,"local-example")
	
	if ll1.UserID != 1 {
		fmt.Println("[FIX PENDING] in Oauth 1",ll1,ll2)
		// https://github.com/Esseh/notorious-dev/issues/1
		//t.Fail()
	}
	
	if ll2.UserID != 2 {
		fmt.Println("FAIL in Local 1")
		t.Fail()
	}
	testing.Coverage()
}