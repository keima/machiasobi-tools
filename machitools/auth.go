package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"
	"appengine"
	"appengine/user"
	"net/http"
)

func Login(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	} else {
		w.Header().Set("Location", "http://localhost:3000")
		w.WriteHeader(http.StatusFound)
	}
}

func Logout(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u != nil {
		url, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	} else {
		w.Header().Set("Location", "http://localhost:3000")
		w.WriteHeader(http.StatusFound)
	}
}

func CheckStatus(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)

	if u == nil {
		// Not logged-in
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
		w.WriteJson(u)
	}
}
