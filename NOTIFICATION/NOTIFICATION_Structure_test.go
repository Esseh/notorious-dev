package NOTIFICATION
import(
	"github.com/Esseh/retrievable"
	"testing"
)

func TestStruct(t*testing.T){
	f := func(retrievable.Retrievable){}
	f(&Notifications{})
}