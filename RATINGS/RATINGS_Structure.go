package RATINGS

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var RatingTable = "Ratings"

type Rating struct{}

func (r *Rating)Key(ctx context.Context,key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, RatingTable, "", key.(int64), nil)	
}