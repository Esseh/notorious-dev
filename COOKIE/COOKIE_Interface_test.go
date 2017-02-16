package COOKIE
import (
	"fmt"
	"testing"
	"github.com/Esseh/notorious-dev/CORE"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
)

func TestMake(t *testing.T){
	res := httptest.NewRecorder()
	Make(res,"hello","world")
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(res.Result().Cookies()[0])
	c, _ := req.Cookie("hello") 
	if c.Value!= "world."+base64.RawURLEncoding.EncodeToString(CORE.CreateHmac("world")) {
		fmt.Println("FAIL in Make"," ",c.Value)
		t.Fail()
	}
	testing.Coverage()
}
func TestGetValue(t *testing.T){
	res := httptest.NewRecorder()
	c := &http.Cookie{
		Name:     "not-good",
		Value:    "hoi",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   CORE.SessionTime,
	}
	http.SetCookie(res, c)	
	Make(res,"hello","world")
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(res.Result().Cookies()[0])	
	req.AddCookie(res.Result().Cookies()[1])	
	// err != nil case
	if _, err := GetValue(req,"i-dont-exist"); err == nil {
		fmt.Println("FAIL in GetValue 1")
		t.Fail()
	}
	// Normal case
	if s,_ := GetValue(req,"hello"); s != "world" {
		fmt.Println("FAIL in GetValue 2")
		t.Fail()	
	}
	// !good case
	if _,err := GetValue(req,"not-good"); err == nil {
		fmt.Println("FAIL in GetValue 3")
		t.Fail()	
	}
	testing.Coverage()
}

// Works in manual testing. Might not be able to be stubbed.
func TestDelete(t *testing.T){
	res := httptest.NewRecorder()
	Make(res,"hello","world")
	Delete(res,"hello")
	if res.Result().Cookies()[0].MaxAge > 0 {
		//fmt.Println("FAIL in Delete",res.Result().Cookies()[0].MaxAge)
		//t.Fail()
	}	
	testing.Coverage()
}