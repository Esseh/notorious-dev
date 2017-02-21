package NOTES
import (
	"errors"
	"strconv"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"google.golang.org/appengine/datastore"
)

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
	validated := VerifyNotePermission(ctx, &Note); if !validated { return errors.New("Permission Error: Not Allowed") }
	if Note.OwnerID == int64(ctx.User.IntID) { Note.Protected = UpdatedNote.Protected }	
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

// Verifies that the currently logged in user is allowed to interact with the Note.
func VerifyNotePermission(ctx CONTEXT.Context, note *Note) bool {
	redirect := strconv.FormatInt(note.OwnerID, 10)
	if note.OwnerID != int64(ctx.User.IntID) && note.Protected {
		ctx.Redirect("/view/"+redirect)
		return false
	}
	return true
}