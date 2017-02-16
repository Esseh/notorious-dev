package CORE
import (
	"strings"
	"net/http"
	"golang.org/x/net/context"
	"github.com/Esseh/notorious-dev/USERS"
	"google.golang.org/appengine"
)
// Constructs an instance of Context
func NewContext(res http.ResponseWriter, req *http.Request) Context{
	user, err := AUTH_GetUserFromSession(req)
	ctx := Context { 
		req: req,
		res: res,
		user: user,
		userException: err,
	}
	ctx.Context = appengine.NewContext(req)
	return ctx
}

// A black box that automatially keeps track of transaction timing for the database
// and stores useful metadata.
type Context struct {
	// The active request 
	req *http.Request
	// The output writer to the user's browser.
	res http.ResponseWriter
	// The currently logged in user.
	user *USERS.User
	// Any problems that occured while logging in.
	userException error
	// transaction timing information.
	context.Context
}

// Returns true if the user is not logged in.
func (ctx Context)AssertLoggedInFailed() bool {
	if ctx.userException != nil {
		path := strings.Replace(ctx.req.URL.Path[1:], "%2f", "/", -1)
		http.Redirect(ctx.res, ctx.req, PATH_AUTH_Login+"?redirect="+path, http.StatusSeeOther)
		return true
	}
	return false
}

// Simplified redirect, useful for general redirects. If the redirect demands a more severe status code use tradition http.Redirect.
func (ctx Context)Redirect(uri string){ http.Redirect(ctx.res, ctx.req, uri, http.StatusSeeOther) }


/// Prints an error page to response and returns a boolean representation of the function executing.
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use if there is no constructive alternative.
func (ctx context)ErrorPage(ErrorTitle string, e error, errCode int) bool {
	if e != nil {
		log.Errorf(ctx, "%s ---- %v\n", ErrorTitle, e)
		if ctx.user == nil {
			ctx.user = &USER_User{}
		}
		args := &struct {
			Header    HeaderData
			ErrorName string
			ErrorDump error
			ErrorCode int
		}{
			HeaderData{ctx, ctx.user, ""}, ErrorTitle, e, errCode,
		}
		ctx.res.WriteHeader(errCode)
		ServeTemplateWithParams(ctx.res, "site-error", args)
		return true
	}
	return false
}

/// Returns to GET responding with FormValue("ErrorResponse")
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use in POST calls accessed from a GET of the same handle.
func (ctx context)BackWithError(err error, errorString string) bool {
	if err != nil {
		path := strings.Replace(ctx.req.URL.Path, "%2f", "/", -1)
		path += "?"+url.QueryEscape("ErrorResponse")+"="+url.QueryEscape(errorString)
		if ctx.req.FormValue("redirect") != "" {
			path += "&"+url.QueryEscape("redirect")+"="+ctx.req.FormValue("redirect")		
		}
		http.Redirect(ctx.res, ctx.req, path, http.StatusSeeOther)
		return true
	}
	return false
}

