package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	PATH_NOTES_New      = "/new"
	PATH_NOTES_View     = "/view/:ID"
	PATH_NOTES_Editor   = "/edit/:ID"
	PATH_NOTES_Edit     = "/edit/"
)

func INIT_NOTES_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_NOTES_New, NOTES_GET_New)
	r.POST(PATH_NOTES_New, NOTES_POST_New)
	r.GET(PATH_NOTES_View, NOTES_GET_View)
	r.GET(PATH_NOTES_Editor, NOTES_GET_Editor)
	r.POST(PATH_NOTES_Edit, NOTES_POST_Editor)
}


func NOTES_GET_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	_,validated := MustLogin(res, req); if !validated { return }
	ServeTemplateWithParams(res, "new-note", MakeHeader(res, req, false, true))
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req); if !validated { return }
	ctx := NewContext(res,req)

	protected, boolConversionError := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", boolConversionError, http.StatusSeeOther) { return }

	
	_, noteKey, err := CreateNewNote(ctx,
		Content{
			Title:   req.FormValue("title"),
			Content: req.FormValue("note"),
		},
		Note{
			OwnerID:   int64(user.IntID),
			Protected: protected,
		},
	)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) { return }
	
	http.Redirect(res, req, "/view/"+strconv.FormatInt(noteKey.IntID(), 10), http.StatusSeeOther)
}

func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req);  if !validated { return }
	ctx := NewContext(res,req)

	ViewNote, ViewContent, err := GetExistingNote(params.ByName("ID"), ctx)
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", err, http.StatusSeeOther) { return }

	owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) { return }

	NoteBody := template.HTML(EscapeString(ViewContent.Content))

	ServeTemplateWithParams(res, "viewNote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
		User, Owner                                *User
	}{
		HeaderData:    *MakeHeader(res, req, false, true),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       NoteBody,
		User:          user,
		Owner:         owner,
	})

}

func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req); if !validated { return }
	ctx := NewContext(res,req)
	ViewNote, ViewContent, err := GetExistingNote(params.ByName("ID"), ctx)
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", err, http.StatusSeeOther) { return }

	validated = VerifyNotePermission(res, req, user, ViewNote); if !validated { return }

	Body := template.HTML(ViewContent.Content)
	ServeTemplateWithParams(res, "editnote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
	}{
		HeaderData:    *MakeHeader(res, req, false, true),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       Body,
	})
}

func NOTES_POST_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req); if !validated { return }
	ctx := NewContext(res,req)

	protbool, boolConversionError := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", boolConversionError, http.StatusSeeOther) { return }

	err := UpdateNoteContent(ctx, res,req ,user ,req.FormValue("notekey"),
		Content{
			Content: EscapeString(req.FormValue("note")),
			Title: req.FormValue("title"),
		},
		Note{
			Protected: protbool,
		},
	)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) { return }	
	http.Redirect(res, req, "/view/"+req.FormValue("notekey"), http.StatusSeeOther)
}