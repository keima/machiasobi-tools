package machitools

import (
	"net/http"
	"machitools/customer"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"github.com/knightso/base/errors"
)

const (
	machiAppURL = "http://machi.p-side.net/"
)

// TODO: 要リファクタリング

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
		if err := createUserIfNeed(c, u); err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		redirect := "/"
		if r.FormValue("redirectTo") == "app" {
			redirect = machiAppURL
		}

		w.Header().Set("Location", redirect)
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
		redirect := "/"
		if r.FormValue("redirectTo") == "app" {
			redirect = machiAppURL
		}

		w.Header().Set("Location", redirect)
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

func createUserIfNeed(c context.Context, u *user.User) error {
	customer := customer.CustomerItem{}

	customer.Init(u)
	if err := customer.Load(c); err != nil {
		if be, ok := err.(*errors.BaseError); ok {
			if be.Cause() == datastore.ErrNoSuchEntity {
				log.Infof(c, "User created: %s", customer.ID)
				customer.Save(c)
			} else {
				return be
			}
		} else {
			return err
		}
	}
	return nil
}
