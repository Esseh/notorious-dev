package NOTES
import(
	"google.golang.org/appengine/aetest"
	"testing"
	"fmt"
)

func TestNote(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestNote")
		panic(1)
	}
	if (&Note{}).Key(ctx,int64(1)) == nil {
		fmt.Println("FAIL Note")
		t.Fail()
	}
}
func TestContent(t*testing.T){
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		fmt.Println("PANIC in TestContent")
		panic(1)
	}
	if (&Content{}).Key(ctx,int64(1)) == nil {
		fmt.Println("FAIL Content")		
		t.Fail()
	}
}