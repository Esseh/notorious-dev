package FOLDERS

import (
	"strconv"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/NOTES"
)

func NewFolder(ctx CONTEXT.Context) string {
	parentFolder := Folder{}
	parentID := ctx.Req.FormValue("ParentID")
	folderName := ctx.Req.FormValue("FolderName")
	
	// Make sure the filepath isn't too long.
	if len(parentID+"/"+folderName) > 400 { return `{"success":false,"code":2}` }
	// Make sure the parent actually exists. If it does then retrieve it.
	if retrievable.GetEntity(ctx,parentID,&parentFolder) != nil { return `{"success":false,"code":1}` }
	// If that folder already exists, just abort now.
	for _,v := range parentFolder.ChildFolders { if v == folderName { return `{"success":false,"code":0}` } }
	// If the parent folder isn't owned by the user accessing it then stop.
	if int64(ctx.User.IntID) != parentFolder.OwnerID { return `{"success":false,"code":3}` }
	// Place the Entry
	child := Folder{
		OwnerID:int64(ctx.User.IntID),
		ParentFolder:parentID,
	}
	_, randomDatabaseErr1 := retrievable.PlaceEntity(ctx,parentID+"/"+folderName,&child)
	if randomDatabaseErr1 != nil { return `{"success":false,"code":1}` }
	// Update Parent
	parentFolder.ChildFolders = append(parentFolder.ChildFolders,folderName)
	_ , randomDatabaseErr2 := retrievable.PlaceEntity(ctx,parentID,&parentFolder)
	if randomDatabaseErr2 != nil { return `{"success":false,"code":1}` }
	// Success
	return `{"success":true,"code":-1}`
}

func DeleteFolder(ctx CONTEXT.Context) string {
	parentFolder := Folder{}
	currentFolder:= Folder{}
	parentID := ctx.Req.FormValue("ParentID")
	folderName := ctx.Req.FormValue("FolderName")
	
	// Not allowed to delete root.
	if parentID == "" { return `{"success":false,"code":4}` }

	// Retrieve the parent and current folders.
	getErr1 := retrievable.GetEntity(ctx,parentID,&parentFolder)
	getErr2 := retrievable.GetEntity(ctx,parentID+"/"+folderName,&currentFolder)
	if getErr1 != nil || getErr2 != nil { return `{"success":false,"code":1}` }

	// Make sure the user owns the current folders.
	if int64(ctx.User.IntID) != parentFolder.OwnerID || int64(ctx.User.IntID) != currentFolder.OwnerID { return `{"success":false,"code":3}` }
	
	// Make sure that the parent and child match.
	if currentFolder.ParentFolder != parentID { return `{"success":false,"code":0}` }
	
	// Detatch currentFolder from parent
	for i,v := range parentFolder.ChildFolders {
		if v == folderName {
			parentFolder.ChildFolders = append(parentFolder.ChildFolders[:i], parentFolder.ChildFolders[i+1:]...)
			break
		}
	}
	// Update Parent
	_, putErr1 := retrievable.PlaceEntity(ctx,parentID,&parentFolder)
	if putErr1 != nil { return `{"success":false,"code":1}` }

	// Perform recursive deletion of child folders.(This will probably need to be additionally supported with a CRON Job)
	var deleteAll func(CONTEXT.Context,string)
	deleteAll = func(ctxi CONTEXT.Context, keyi string){
		targetFolder := Folder{}
		retrievable.GetEntity(ctxi,keyi,&targetFolder)
		for _,v := range targetFolder.ChildFolders {
			deleteAll(ctxi, keyi + "/" + v)
		}
		retrievable.DeleteEntity(ctxi,(&Folder{}).Key(ctxi,keyi))
	}
	deleteAll(ctx,parentID+"/"+folderName)
	
	// Return Success
	return `{"success":true,"code":-1}`
}

func OpenFolder(ctx CONTEXT.Context) string {
	validFolder := ""
	validNotes  := ""
	var updateFolder bool
	successfullyOpenedNotes := make([]int64,0)
	folder := Folder{}
	folderID := ctx.Req.FormValue("FolderID")

	// Open the Folder
	if retrievable.GetEntity(ctx,folderID,&folder) != nil { return `{"success":false, "code":1}`}
	
	// Retrieve each Folder
	for _,v := range folder.ChildFolders { validFolder += `"`+v+`"`+"," }
	// Trim last comma
	if validFolder != "" { validFolder = string(validFolder[:len(validFolder)-1]) }
	
	// Retrieve each Note
	for _,v := range folder.ChildNotes {
		// Get the note, the content of the note, and check if it actually exists.
		note, content, openErr :=  NOTES.GetExistingNote(ctx,strconv.FormatInt(v,10))
		// If it doesn't exist then mark for updating.
		if openErr != nil { updateFolder = true; continue }
		// If the note was succesfully opened, make a note of it.
		successfullyOpenedNotes = append(successfullyOpenedNotes,v)
		// Only consider the notes that the user is allowed to see.
		if NOTES.CanViewNote(note,ctx.User){
			if NOTES.CanEditNote(note,ctx.User){
				validNotes += `{"title":"`+content.Title+`","id":`+strconv.FormatInt(v,10)+`,"noEdit":false},`
			} else {
				validNotes += `{"title":"`+content.Title+`","id":`+strconv.FormatInt(v,10)+`,"noEdit":true},`
			}
		}
	}
	if validNotes != "" { validNotes = string(validNotes[:len(validNotes)-1]) }
	
	// If the folder needs updating
	if updateFolder {
		folder.ChildNotes = successfullyOpenedNotes
		retrievable.PlaceEntity(ctx,folderID,&folder)
	}
	
	// Return Success
	return `{"success":true,"code":-1,"folders":[`+validFolder+`],"notes":[`+validNotes+`]}`
}

func AddNote(ctx CONTEXT.Context) string {
	folder := Folder{}
	note := NOTES.Note{}
	folderID := ctx.Req.FormValue("FolderID")
	noteID, _ := strconv.ParseInt(ctx.Req.FormValue("NoteID"),10,64)
	// Get folder and note
	getErr1 := retrievable.GetEntity(ctx,folderID,&folder)
	getErr2 := retrievable.GetEntity(ctx,noteID,&note)
	if getErr1 != nil || getErr2 != nil { return `{"success":false,"code":1}` }
	// If the parent folder isn't owned by the user accessing it then stop.
	if int64(ctx.User.IntID) != folder.OwnerID { return `{"success":false,"code":3}` }
	// Make sure the reference isn't already inside.
	for _,v := range folder.ChildNotes { if v == noteID { return `{"success":false,"code":0}`} }
	// Attatch the reference and Update
	folder.ChildNotes = append(folder.ChildNotes,noteID)
	_ , placeErr := retrievable.PlaceEntity(ctx,folderID,&folder)
	if placeErr != nil { return `{"success":false,"code":1}` }
	return `{"success":true,"code":-1}`
}

func RemoveNote(ctx CONTEXT.Context) string {
	folder := Folder{}
	note := NOTES.Note{}
	folderID := ctx.Req.FormValue("FolderID")
	noteID, _ := strconv.ParseInt(ctx.Req.FormValue("NoteID"),10,64)
	// Get folder and note
	getErr1 := retrievable.GetEntity(ctx,folderID,&folder)
	getErr2 := retrievable.GetEntity(ctx,noteID,&note)
	if getErr1 != nil || getErr2 != nil { return `{"success":false,"code":1}` }
	// If the parent folder isn't owned by the user accessing it then stop.
	if int64(ctx.User.IntID) != folder.OwnerID { return `{"success":false,"code":3}` }
	// Build References Excluding the One Removed
	newChildNotes := make([]int64,0)
	for _,v := range folder.ChildNotes { if v != noteID { newChildNotes = append(newChildNotes,v) } }
	folder.ChildNotes = newChildNotes
	_ , placeErr := retrievable.PlaceEntity(ctx,folderID,&folder)
	if placeErr != nil { return `{"success":false,"code":1}` }
	return `{"success":true,"code":-1}`
}

func InitializeRoot(ctx CONTEXT.Context) string {
	folderID := ctx.Req.FormValue("RootID")
	userID, _ := strconv.ParseInt(folderID,10,64)

	// Invalid User?
	if retrievable.GetEntity(ctx,userID,&USERS.User{}) != nil { return `{"success":false,"code":5}` }
	
	// Already Exists?
	folder := Folder{}
	if retrievable.GetEntity(ctx,folderID,&folder) == nil { return `{"success":false,"code":6}` }
	// Create Root
	_ , placeErr := retrievable.PlaceEntity(ctx,folderID,&Folder{
		IsRoot:true,
		OwnerID:userID,
	})
	if placeErr != nil { return `{"success":true,"code":1}` }
	// Success
	return `{"success":true,"code":-1}`
}