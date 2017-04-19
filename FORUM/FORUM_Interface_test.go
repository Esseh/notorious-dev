package FORUM
import (
	"fmt"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/aetest"
	"testing"
)

func TestMakeChain(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in RegisterAdminAPI")
		panic(1)
	}
	ch := CategoryHeader{}
	ca := Category{}
	fo := Forum{}
	th := Thread{}	
	po := Post{}
	
	MakeCategory(ctx, "testCategory")
	retrievable.GetEntity(ctx,nil,&ch)
	MakeForum(ctx,"testForum","testForumDescription",ch.Categories[0])
	retrievable.GetEntity(ctx,ch.Categories[0],&ca)
	MakeThread(ctx,"testTitle",ca.Forums[0])
	retrievable.GetEntity(ctx,ca.Forums[0],&fo)
	MakePost(ctx,"testPoster","postBody",fo.Threads[0])
	retrievable.GetEntity(ctx,fo.Threads[0],&th)
	retrievable.GetEntity(ctx,th.Posts[0],&po)
	
	if po.Poster != "testPoster" {
		fmt.Println("FAIL MakeChain 1")
		t.Fail()
	}
	if th.Title != "testTitle"   {
		fmt.Println("FAIL MakeChain 2")	
		t.Fail()
	}
	if fo.Title != "testForum"   {
		fmt.Println("FAIL MakeChain 3")	
		t.Fail()
	}
	if ca.Name != "testCategory" {
		fmt.Println("FAIL MakeChain 4")	
		t.Fail()
	}
}