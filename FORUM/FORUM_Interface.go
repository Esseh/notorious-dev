package FORUM
import(
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
)

func MakeCategory(ctx context.Context, name string){
	ch := CategoryHeader{}
	retrievable.GetEntity(ctx,nil,&ch)
	key , _ := retrievable.PlaceEntity(ctx,int64(0),&Category{Name:name,})
	ch.Categories = append([]int64{key.IntID()},ch.Categories...)
	retrievable.PlaceEntity(ctx,nil,&ch)
}

func MakeForum(ctx context.Context,name string,description string,categoryid int64){
	c := Category{}
	retrievable.GetEntity(ctx,categoryid,&c)
	key , _ := retrievable.PlaceEntity(ctx,int64(0),&Forum{Title:name,Description:description,})
	c.Forums = append([]int64{key.IntID()},c.Forums...)
	retrievable.PlaceEntity(ctx,categoryid,&c)
}

func MakeThread(ctx context.Context,name string,forumid int64){
	c := Forum{}
	retrievable.GetEntity(ctx,forumid,&c)
	key , _ := retrievable.PlaceEntity(ctx,int64(0),&Thread{Title:name})
	c.Threads = append([]int64{key.IntID()},c.Threads...)
	retrievable.PlaceEntity(ctx,forumid,&c)
}

func MakePost(ctx context.Context,poster string,body string,threadid int64){
	c := Thread{}
	retrievable.GetEntity(ctx,threadid,&c)
	key , _ := retrievable.PlaceEntity(ctx,int64(0),&Post{Poster:poster,Body:body,})
	c.Posts = append([]int64{key.IntID()},c.Posts...)
	retrievable.PlaceEntity(ctx,threadid,&c)
}