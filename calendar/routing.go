package calendar

import (
	"net/http"

	"github.com/keima/machiasobi-tools/customer"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func GetFavList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || u.ID == "" {
		rest.Error(w, "Login required.", http.StatusUnauthorized)
		return
	}

	key := datastore.NewKey(c, customer.KindName, u.ID, 0, nil)

	itemList, err := LoadAll(c, key)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&itemList)
}

func PostFav(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || u.ID == "" {
		rest.Error(w, "Login required.", http.StatusUnauthorized)
		return
	}

	key := datastore.NewKey(c, customer.KindName, u.ID, 0, nil)

	item := FavItem{
		CalendarID: r.PathParam("calId"),
		EventID:    r.PathParam("eventId"),
	}

	if err := item.Save(c, key); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func DeleteFav(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || u.ID == "" {
		rest.Error(w, "Login required.", http.StatusUnauthorized)
		return
	}

	key := datastore.NewKey(c, customer.KindName, u.ID, 0, nil)

	item := FavItem{
		CalendarID: r.PathParam("calId"),
		EventID:    r.PathParam("eventId"),
	}

	if err := item.Delete(c, key); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
