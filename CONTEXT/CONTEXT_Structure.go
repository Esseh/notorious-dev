package CONTEXT
import (
	"strings"
	"net/http"
	"net/url"
	"golang.org/x/net/context"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/USERS"
	"google.golang.org/appengine"
)

// Header Data,
// Present in most template executions. (Unless it's an internal it should be assumed to be used.)
// Placed here to minimize cross-dependencies.
type HeaderData struct {
	Ctx          context.Context
	User         *USERS.User
	CurrentPath  string
}

// Constructs an instance of Context
func NewContext(res http.ResponseWriter, req *http.Request) Context{
	ctx := Context { 
		Req: req,
		Res: res,
	}
	ctx.Context = appengine.NewContext(req)
	ctx.User, ctx.UserException = USERS.GetUserFromSession(ctx.Context,req)
	return ctx
}

// A black box that automatially keeps track of transaction timing for the database
// and stores useful metadata.
type Context struct {
	// The active request 
	Req *http.Request
	// The output writer to the user's browser.
	Res http.ResponseWriter
	// The currently logged in user.
	User *USERS.User
	// Any problems that occured while logging in.
	UserException error
	// transaction timing information.
	context.Context
}

// Returns true if the user is not logged in.
func (ctx Context)AssertLoggedInFailed() bool {
	if ctx.UserException != nil {
		path := strings.Replace(ctx.Req.URL.Path[1:], "%2f", "/", -1)
		http.Redirect(ctx.Res, ctx.Req, PATHS.AUTH_Login+"?redirect="+path, http.StatusSeeOther)
		return true
	}
	return false
}

// Simplified redirect, useful for general redirects. If the redirect demands a more severe status code use tradition http.Redirect.
func (ctx Context)Redirect(uri string){ http.Redirect(ctx.Res, ctx.Req, uri, http.StatusSeeOther) }


/// Prints an error page to response and returns a boolean representation of the function executing.
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use if there is no constructive alternative.
func (ctx Context)ErrorPage(ErrorTitle string, e error, errCode int) bool {
	if e != nil {
		if ctx.User == nil {
			ctx.User = &USERS.User{}
		}
		args := &struct {
			Header    HeaderData
			ErrorName string
			ErrorDump error
			ErrorCode int
		}{
			HeaderData{ctx, ctx.User, ""}, ErrorTitle, e, errCode,
		}
		ctx.Res.WriteHeader(errCode)
		CORE.ServeTemplateWithParams(ctx.Res, "site-error", args)
		return true
	}
	return false
}

/// Returns to GET responding with FormValue("ErrorResponse")
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use in POST calls accessed from a GET of the same handle.
func (ctx Context)BackWithError(err error, errorString string) bool {
	if err != nil {
		path := strings.Replace(ctx.Req.URL.Path, "%2f", "/", -1)
		path += "?"+url.QueryEscape("ErrorResponse")+"="+url.QueryEscape(errorString)
		if ctx.Req.FormValue("redirect") != "" {
			path += "&"+url.QueryEscape("redirect")+"="+ctx.Req.FormValue("redirect")		
		}
		http.Redirect(ctx.Res, ctx.Req, path, http.StatusSeeOther)
		return true
	}
	return false
}

