package event

import (
	"net/http"
	"strconv"
	"time"

	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/user"
)

func GetEventList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	first, err := strconv.Atoi(r.FormValue("first"))
	if err != nil {
		first = 0
	}

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		size = 10
	}

	startAt, err := parseDate(r.FormValue("startAt"))
	if err != nil {
		startAt = time.Time{}
	}

	endAt, err := parseDate(r.FormValue("endAt"))
	if err != nil {
		endAt = time.Time{}
	}

	publicOnly := true
	if r.Request.FormValue("private") == "true" {
		if u := user.Current(c); u != nil && user.IsAdmin(c) {
			publicOnly = false
		}
	}

	items, err := LoadAll(c, first, size, startAt, endAt, publicOnly)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&items)
}

func GetEvent(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	item := EventItem{}

	err := item.Load(c, r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func PostEvent(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := EventItem{}

	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.Author = u.String()

	var intId int64
	var err error

	if id := r.PathParam("id"); id != "" {
		intId, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if intId == 0 {
		err = item.Save(c)
	} else {
		err = item.SaveUpdate(c, intId)
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// 2014-10-11T15:00:00.000Z のような形式を処理する
func parseDate(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

/*
func GetEventList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

}
*/
