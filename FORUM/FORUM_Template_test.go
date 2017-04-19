package FORUM
import (
	"fmt"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/aetest"
	"testing"
)

func TestIsAdmin(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in IsAdmin")
		panic(1)
	}
	retrievable.PlaceEntity(ctx,nil,&AdminHeader{UserIDS:[]int64{1},})
	if IsAdmin(ctx,2)  {
		t.Fail()
		fmt.Println("FAIL IsAdmin 1")
	}
	if !IsAdmin(ctx,1) {
		t.Fail()
		fmt.Println("FAIL IsAdmin 2")
	}
}

func TestGetCategories(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetCategories")
		panic(1)
	}
	retrievable.PlaceEntity(ctx,nil,&CategoryHeader{Categories:[]int64{1,2},})	
	retrievable.PlaceEntity(ctx,int64(1),&Category{Name:"1",})	
	retrievable.PlaceEntity(ctx,int64(2),&Category{Name:"2",})

	out := GetCategories(ctx)
	if out[0].Name != "1" || out[1].Name != "2" {
		fmt.Println("FAIL GetCategories")
		t.Fail()
	}
}

func TestGetForums(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetForums")
		panic(1)
	}
	retrievable.PlaceEntity(ctx,int64(1),&Forum{Title:"1",})	
	retrievable.PlaceEntity(ctx,int64(2),&Forum{Title:"2",})
	
	out := GetForums(ctx,Category{Forums:[]int64{1,2},})
	if out[0].Title != "1" || out[1].Title != "2" {
		fmt.Println("FAIL GetForums")
		t.Fail()
	}	
}

func TestGetThreads(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetThreads")
		panic(1)
	}
	retrievable.PlaceEntity(ctx,int64(1),&Thread{Title:"1",})	
	retrievable.PlaceEntity(ctx,int64(2),&Thread{Title:"2",})
	out := GetThreads(ctx,Forum{Threads:[]int64{1,2},})
	if out[0].Title != "1" || out[1].Title != "2" {
		fmt.Println("FAIL GetThreads")
		t.Fail()
	}	
}

func TestGetPosts(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestGetPosts")
		panic(1)
	}
	retrievable.PlaceEntity(ctx,int64(1),&Post{Poster:"1",})	
	retrievable.PlaceEntity(ctx,int64(2),&Post{Poster:"2",})	
	out := GetPosts(ctx,Thread{Posts:[]int64{1,2},})
	if out[0].Poster != "1" || out[1].Poster != "2" {
		fmt.Println("FAIL GetPosts")
		t.Fail()
	}	
}