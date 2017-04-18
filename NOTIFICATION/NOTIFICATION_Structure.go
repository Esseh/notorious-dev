package NOTIFICATION

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	NotificationTable = "Notifications"
)

type Notifications struct {
	NotificationsPending int64
	Notifications []string
}


func (n *Notifications) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, NotificationTable, "", key.(int64), nil)
}