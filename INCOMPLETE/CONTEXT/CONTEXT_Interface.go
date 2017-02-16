/// Prints an error page to response and returns a boolean representation of the function executing.
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use if there is no constructive alternative.
func ERROR_Page(ctx Context, ErrorTitle string, e error, errCode int) bool {
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
func ERROR_Back(ctx Context, err error, errorString string) bool {
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

