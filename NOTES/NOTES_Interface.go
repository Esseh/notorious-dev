package NOTES
import (
	"errors"
	"strconv"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"google.golang.org/appengine/datastore"
)

func API_SaveCopy(ctx CONTEXT.Context) string { 
	OriginalNote    := Note{}
	OriginalContent := Content{}

	// Parse Key. Exit if bad Key
	noteID, err := strconv.ParseInt(ctx.Req.FormValue("NoteID"), 10, 64)
	if err != nil { return `{"success":false}` }
	
	// Get Note. Exit if doesn't exist.
	err = retrievable.GetEntity(ctx,noteID,&OriginalNote)
	if err != nil { return `{"success":false}` }
	
	// Ensure that the logged in user is actually allowed to look at the note.
	if !CanViewNote(&OriginalNote,ctx.User) { return `{"success":false}` }
	
	// Get Original Content, check for random failure
	err = retrievable.GetEntity(ctx,OriginalNote.ContentID,&OriginalContent)
	if err != nil { return `{"success":false}` }
	
	// Make Copy of Original Content and attatch new Note Header to it owned by the current user.
	contentKey , err := retrievable.PlaceEntity(ctx,int64(0),&OriginalContent)
	if err != nil { return `{"success":false}` }
	noteKey , err := retrievable.PlaceEntity(ctx,int64(0),&Note{OwnerID:int64(ctx.User.IntID),ContentID:contentKey.IntID(),})
	if err != nil { return `{"success":false}` }
	
	// Success
	return `{"success":true,"CopyID":`+strconv.FormatInt(noteKey.IntID(),10)+`}`
}

// Given a Content and Note it will construct instances of each, tie them together in the database and provide their keys.
func CreateNewNote(ctx CONTEXT.Context,NewContent Content,NewNote Note) (*datastore.Key,*datastore.Key,error) {
	contentKey, _ := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	NewNote.ContentID = contentKey.IntID()
	noteKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	return contentKey,noteKey,err
}

// Updates a note and its content based on the given id.
func UpdateNoteContent(ctx CONTEXT.Context,id string,UpdatedContent Content, UpdatedNote Note) error {
	Note := Note{}
	noteID, err := strconv.ParseInt(id, 10, 64); if err != nil { return err }
	err = retrievable.GetEntity(ctx, noteID, &Note); if err != nil { return err }
	validated := CanEditNote(&Note,ctx.User); if !validated { return errors.New("Permission Error: Not Allowed") }
	if Note.OwnerID == int64(ctx.User.IntID) { 
		Note.PublicallyViewable = UpdatedNote.PublicallyViewable
		Note.PublicallyEditable = UpdatedNote.PublicallyEditable
		Note.Collaborators = UpdatedNote.Collaborators
	}	
	_, err = retrievable.PlaceEntity(ctx, noteID, &Note); if err != nil { return err }
	_, err = retrievable.PlaceEntity(ctx, Note.ContentID, &UpdatedContent); return err
}

// Retrieves an existing note and it's content by it's id.
func GetExistingNote(ctx CONTEXT.Context,id string)(*Note,*Content,error){
	RetrievedNote := &Note{}
	RetrievedContent := &Content{}
	NoteKey, _ := strconv.ParseInt(id, 10, 64)
	retrievable.GetEntity(ctx, NoteKey, RetrievedNote)
	err := retrievable.GetEntity(ctx, RetrievedNote.ContentID, RetrievedContent)
	return RetrievedNote,RetrievedContent,err
}

// Gets all the notes made by the AUTH_User
func GetAllNotes(ctx CONTEXT.Context, userID int64) ([]NoteOutput, error) {
	q := datastore.NewQuery(NoteTable).Filter("OwnerID =", userID)
	res := []Note{}
	output := make([]NoteOutput, 0)
	keys, _ := q.GetAll(ctx, &res)
	for i, _ := range res {
		var c Content
		err := retrievable.GetEntity(ctx, res[i].ContentID, &c)
		if err != nil {
			return nil, err
		}
		output = append(output, NoteOutput{keys[i].IntID(), res[i], c})
	}
	return output, nil
}


func CanViewNote(note *Note,user *USERS.User) bool {
	uid := int64(user.IntID)
	if uid == note.OwnerID { return true }
	for _, v := range note.Collaborators {
		if uid == v {
			return true
		}
	}
	return note.PublicallyViewable
}

func CanEditNote(note *Note,user *USERS.User) bool {
	uid := int64(user.IntID)
	if uid == note.OwnerID { return true }
	for _, v := range note.Collaborators {
		if uid == v {
			return true
		}
	}	
	return note.PublicallyEditable
}