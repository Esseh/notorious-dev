package FORUM

import(
	"github.com/Esseh/retrievable"
	"testing"
)

func TestStruct(t*testing.T){
	f := func(retrievable.Retrievable){}
	f(&AdminHeader{})
	f(&CategoryHeader{})
	f(&Category{})
	f(&Forum{})
	f(&Thread{})
	f(&Post{})
}