package FORUM
import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type (
	AdminHeader struct {
		UserIDS []int64
	}
	CategoryHeader struct {
		Categories []int64
	}
	Category struct {
		Name string
		Forums []int64
	}
	Forum struct {
		Title string
		Description string
		Threads []int64
	}
	Thread struct {
		Title string
		Posts []int64
	}
	Post struct {
		Poster string
		Body string
	}
)

func (a *AdminHeader) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "AdminHeader", "", 1, nil)
}

func (a *CategoryHeader) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "CategoryHeader", "", 1, nil)
}

func (a *Category) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Categories", "", k.(int64), nil)
}
func (a *Forum) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Forums", "", k.(int64), nil)
}
func (a *Thread) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Threads", "", k.(int64), nil)
}
func (a *Post) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Posts", "", k.(int64), nil)
}