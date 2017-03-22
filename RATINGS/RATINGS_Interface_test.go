package RATINGS
import(
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
	"github.com/Esseh/retrievable"
	"net/http/httptest"
	"net/url"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/NOTES"
)

func TestGetRating(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in GetRating")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&NOTES.Note{})
	
	// Make ctx1
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("NoteID","2")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	
	// Make ctx2
	req2 := httptest.NewRequest("GET", "/", nil)
	values2 := url.Values{}; 
	values2.Add("NoteID","1")
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx; 
	ctx2.Req = req2
	
	// Note Doesn't Exist Case
	if GetRating(ctx1) != `{"success":false,"code":1}` {
		fmt.Println("FAIL GetRating 1")		
		t.Fail()
	}
	// Note Exists Case
	if x := GetRating(ctx2); x != `{"success":true,"totalRating":0,"code":-1}` {
		fmt.Println("FAIL GetRating 2",x)
		t.Fail()
	}
	// Assert Success
	if retrievable.GetEntity(ctx,int64(1),&Rating{}) != nil {
		fmt.Println("FAIL GetRating 3")
		t.Fail()
	}
}

func TestSetRating(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in SetRating")
		panic(1)
	}
	// Stub Database
	retrievable.PlaceEntity(ctx,int64(1),&NOTES.Note{})
	retrievable.PlaceEntity(ctx,int64(1),&Rating{})
	retrievable.PlaceEntity(ctx,int64(1),&USERS.User{})
	
	// Make ctx0
	req0 := httptest.NewRequest("GET", "/", nil)
	values0 := url.Values{}; 
	values0.Add("NoteID","1")
	values0.Add("RatingValue","5")
	req0.Form = values0
	ctx0 := CONTEXT.Context{}
	ctx0.Context = ctx; 
	ctx0.Req = req0
	ctx0.User = &USERS.User{}
	
	// Make ctx1
	req1 := httptest.NewRequest("GET", "/", nil)
	values1 := url.Values{}; 
	values1.Add("NoteID","2")
	values1.Add("RatingValue","5")
	req1.Form = values1
	ctx1 := CONTEXT.Context{}
	ctx1.Context = ctx; 
	ctx1.Req = req1
	ctx1.User = &USERS.User{IntID:retrievable.IntID(int64(1)),}
	
	// Make ctx2
	req2 := httptest.NewRequest("GET", "/", nil)
	values2 := url.Values{}; 
	values2.Add("NoteID","1")
	values2.Add("RatingValue","6")
	req2.Form = values2
	ctx2 := CONTEXT.Context{}
	ctx2.Context = ctx; 
	ctx2.Req = req2
	ctx2.User = &USERS.User{IntID:retrievable.IntID(int64(1)),}
	
	// Make ctx3
	req3 := httptest.NewRequest("GET", "/", nil)
	values3 := url.Values{}; 
	values3.Add("NoteID","1")
	values3.Add("RatingValue","3")
	req3.Form = values3
	ctx3 := CONTEXT.Context{}
	ctx3.Context = ctx; 
	ctx3.Req = req3
	ctx3.User = &USERS.User{IntID:retrievable.IntID(int64(1)),}	
	
	// Make ctx4
	req4 := httptest.NewRequest("GET", "/", nil)
	values4 := url.Values{}; 
	values4.Add("NoteID","1")
	values4.Add("RatingValue","4")
	req4.Form = values4
	ctx4 := CONTEXT.Context{}
	ctx4.Context = ctx; 
	ctx4.Req = req4
	ctx4.User = &USERS.User{IntID:retrievable.IntID(int64(2)),}	
	
	// Make ctx5
	req5 := httptest.NewRequest("GET", "/", nil)
	values5 := url.Values{}; 
	values5.Add("NoteID","1")
	req5.Form = values5
	ctx5 := CONTEXT.Context{}
	ctx5.Context = ctx; 
	ctx5.Req = req5
	
	// Not Logged In (Other)
	if SetRating(ctx0) != `{"success":false,"code":0}` {
		fmt.Println("FAIL SetRating 0")
		t.Fail()
	} 
	// Rating Doesn't Exist Case
	if SetRating(ctx1) != `{"success":false,"code":1}` {
		fmt.Println("FAIL SetRating 1")
		t.Fail()
	} 
	// Value Too Large/Small (Other)
	if SetRating(ctx2) != `{"success":false,"code":0}` {
		fmt.Println("FAIL SetRating 2")
		t.Fail()
	} 
	// Success Case
	if SetRating(ctx3) != `{"success":true,"code":-1}` {
		fmt.Println("FAIL SetRating 3")
		t.Fail()
	} 	
	// Assert Success
	if x := GetRating(ctx5); x != `{"success":true,"totalRating":3,"code":-1}` {
		fmt.Println("FAIL SetRating 4",x)
		t.Fail()	
	}
	// Run Additional Mock Rating
	SetRating(ctx4)
	if x := GetRating(ctx5); x != `{"success":true,"totalRating":3.5,"code":-1}` {
		fmt.Println("FAIL SetRating 5",x)
		t.Fail()	
	}
}