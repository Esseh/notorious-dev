package FOLDERS
import(
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Folder struct{
	IsRoot bool
	OwnerID int64
	ParentFolder string
	ChildFolders []string
	ChildNotes []int64
}

var FolderTable = "Folders"

func (f *Folder)Key(ctx context.Context,key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, FolderTable, key.(string), 0, nil)	
}