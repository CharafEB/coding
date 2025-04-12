package fileserver

import "net/http"

type Files struct {
	HanderlFiles interface {
		HandelWeb(w http.ResponseWriter, r *http.Request)
		HandelOutPutCss(w http.ResponseWriter, r *http.Request)
		Handeltest(w http.ResponseWriter, r *http.Request)
		UploadHandler(w http.ResponseWriter, r *http.Request)
	}

	HandelHomePage interface {
		GetHomePage(w http.ResponseWriter, r *http.Request)
		GetWebEvent(w http.ResponseWriter, r *http.Request)
		GetAddSection(w http.ResponseWriter, r *http.Request)
	}
}
