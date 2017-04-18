package NOTES
import(
	"github.com/Esseh/retrievable"
	"testing"
)

func TestStruct(t*testing.T){
	f := func(retrievable.Retrievable){}
	f(&Note{})
	f(&Content{})
	f(&Subscription{})
	f(&SubscriptionHeader{})
}