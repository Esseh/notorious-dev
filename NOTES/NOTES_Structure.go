package NOTES

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	NoteTable = "NotePermissions"
)

// Contains the metadata of a Note and represents the existance of one.
type Note struct {
	// The ID to the AUTH_User owner.
	OwnerID int64
	// View/Edit Permissions for the Public
	PublicallyViewable,PublicallyEditable bool
	// Collaborators who bypass the security check.
	Collaborators []int64
	// The ID to the related Content
	ContentID int64
}

// The actual content referred to by Note
type Content struct {
	// The title and content respectively.
	Title, Content string
}

func (n *Note) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, NoteTable, "", key.(int64), nil)
}
func (n *Content) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "NoteContents", "", key.(int64), nil)
}


// A convience structure used with GetAllNotes, provides a clean managable output on the front end.
type NoteOutput struct {
	// The ID for a note.
	ID      int64
	// The Note Meta Data
	Data    Note
	// The Note Content
	Content Content
}