package RATINGS

import (
	"fmt"
	"testing"
	"google.golang.org/appengine/aetest"
)

func TestRating(t *testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestRating")
		panic(1)
	}
	// No assertion to check, if it compiles it fulfills the interface.
	(&Rating{}).Key(ctx,int64(1))
}