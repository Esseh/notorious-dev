package USERS
import(
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/aetest"
	"testing"
	"fmt"
)

func TestUser(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestUser")
		panic(1)
	}
	if (&User{}).Key(ctx,retrievable.IntID(1)) == nil {
		fmt.Println("FAIL TestUser 1")
		t.Fail()
	}
	if (&User{}).Key(ctx,int64(1)) == nil {
		fmt.Println("FAIL TestUser 2")
		t.Fail()		
	}
	(&User{Email: "test@test.com",}).toEncrypt()
	err = (&User{}).fromEncrypt(&EncryptedUser{Email:"something-not-encrypted"})
	if err == nil {
		fmt.Println("FAIL TestUser 3")
		t.Fail()		
	}
	err  = (&User{}).fromEncrypt((&User{Email:"",}).toEncrypt())
	if err != nil {
		fmt.Println("FAIL TestUser 4")
		t.Fail()	
	}
	u1 := &User{Email:"hoi",}
	u2 := &User{}
	// error != nil case
	err = u1.Unserialize([]byte("{{{{}"))
	if err == nil {
		fmt.Println("FAIL TestUser 5")
		t.Fail()
	}
	marshalledData := u1.Serialize()
	u2.Unserialize(marshalledData)
	if u1.Email != u2.Email {
		fmt.Println("FAIL TestUser 6",u1.Email,u2.Email)
		t.Fail()
	}
	testing.Coverage()
}
func TestSession(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestSession")
		panic(1)
	}
	s := &Session{}
	k := s.Key(ctx,int64(1))
	if k == nil {
		fmt.Println("FAIL TestSession 1")
		t.Fail()	
	}
	s.StoreKey(k)
	if s.ID != int64(1){
		fmt.Println("FAIL TestSession 2")
		t.Fail()
	}
	testing.Coverage()
}