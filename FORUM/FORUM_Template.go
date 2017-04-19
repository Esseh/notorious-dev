package FORUM
import(
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
)

func IsAdmin(ctx context.Context,UserID int64) bool { 
	ah := AdminHeader{}
	retrievable.GetEntity(ctx,nil,&ah)
	for _,v := range ah.UserIDS {
		if v == UserID { return true }
	}
	return false 
}

func GetCategories(ctx context.Context) []Category { 
	ch := CategoryHeader{}
	retrievable.GetEntity(ctx,nil,&ch)
	out := []Category{}
	for _,v := range ch.Categories {
		cat := Category{}
		if retrievable.GetEntity(ctx,v,&cat) == nil {
			out = append(out,cat)
		}
	}
	return out
}

func GetForums(ctx context.Context, c Category) []Forum { 
	out := []Forum{}
	for _,v := range c.Forums {
		f := Forum{}
		if retrievable.GetEntity(ctx,v,&f) == nil {
			out = append(out,f)
		}
	}
	return out
}

func GetThreads(ctx context.Context, f Forum) []Thread { 
	out := []Thread{}
	for _,v := range f.Threads {
		t := Thread{}
		if retrievable.GetEntity(ctx,v,&t) == nil {
			out = append(out,t)
		}
	}
	return out
}

func GetPosts(ctx context.Context, t Thread) []Post { 
	out := []Post{}
	for _,v := range t.Posts {
		p := Post{}
		if retrievable.GetEntity(ctx,v,&p) == nil {
			out = append(out,p)
		}
	}
	return out
}