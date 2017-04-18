package NOTIFICATION
import (
	"golang.org/x/net/context"
	USER_CONTEXT "github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/retrievable"
)

// simple wrapper function, shouldn't require testing.
func GetNotifications(ctx context.Context, UserID int64) *Notifications {
	n := &Notifications{}				// Heap Object
	retrievable.GetEntity(ctx,UserID,n) // Heap Object Passed in Normally (pointer)
	return n							// Retrieved object returned.
}

func ClearNotificationsAPI(ctx USER_CONTEXT.Context) string {
	n := &Notifications{}
	retrievable.GetEntity(ctx,int64(ctx.User.IntID),n)
	n.NotificationsPending = 0
	retrievable.PlaceEntity(ctx,int64(ctx.User.IntID),n)
	return `{"success":true}`
}