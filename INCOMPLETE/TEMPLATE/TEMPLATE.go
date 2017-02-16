package TEMPLATE

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"regexp"

	"github.com/Esseh/retrievable"
	humanize "github.com/dustin/go-humanize" // russross markdown parser
	"golang.org/x/net/context"
)



func init() {
	// Tie functions into template here with ... "functionName":theFunction,
	funcMap := template.FuncMap{
		"getAvatarURL":  getAvatarURL,
		"getUser":       AUTH_GetUserFromID,
		"humanize":      humanize.Time,
		"humanizeSize":  humanize.Bytes,
		"monthfromtime": monthfromtime,
		"yearfromtime":  yearfromtime,
		"dayfromtime":   dayfromtime,
		"findsvg":       FindSVG,
		"findtemplate":  FindTemplate,
		"inc":           Inc,
		"addCtx":        addCtx,
		"getDate":       getDate,
		"toInt":		 toInt,
		// "isOwner":       isOwner,
		"parse": EscapeString,
	} // Load up all templates.
	tpl = template.New("").Funcs(funcMap)
	tpl = template.Must(tpl.ParseGlob("templates/*"))

}

// Constructs the header.
// As the header gets more complex(such as capturing the current path)
// the need for such a helper function increases.
func MakeHeader(ctx Context) *HeaderData {
	oldCookie, err := COOKIE_GetValue(ctx.req, "session")
	if err == nil { COOKIE_Make(ctx.res, "session", oldCookie) }
	redirectURL := ctx.req.URL.Path[1:]
	if redirectURL == "login" || redirectURL == "register" || redirectURL == "elevatedlogin" {
		redirectURL = ctx.req.URL.Query().Get("redirect")
	}
	return &HeaderData{
		ctx, ctx.user, redirectURL,
	}
}

