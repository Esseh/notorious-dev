package PM
import(
	"time"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type PrivateMessageHeader struct{
	Messages []int64
}
type PrivateMessage struct{
	Sender, Receiver, Title, Content string
	DateSent time.Time
}

var PrivateMessageHeaderTable = "PrivateMessageHeaders"
var PrivateMessageTable = "PrivateMessages"

func (f *PrivateMessageHeader)Key(ctx context.Context,key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, PrivateMessageHeaderTable, "", key.(int64), nil)	
}

func (f *PrivateMessage)Key(ctx context.Context,key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, PrivateMessageTable, "", key.(int64), nil)	
}