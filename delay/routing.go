package delay

import (
	"errors"
	"net/http"

	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/user"
)

func GetDelay(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	placeName, err := checkPlaceName(r.PathParam("place"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := LoadLatest(c, placeName)
	if err != nil {
		if err == ErrItemNotFound {
			rest.Error(w, err.Error(), http.StatusNotFound)
		} else {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteJson(item)
}

func PostDelay(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	placeName, err := checkPlaceName(r.PathParam("place"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := DelayItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.Author = u.String()

	if err := item.Save(c, placeName); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func checkPlaceName(placeName string) (string, error) {
	switch placeName {
	case "bizan", "shinmachi", "corne", "cinema_entry", "awagin", "bunka", "ryogoku":
		return placeName, nil
	default:
		return "", errors.New("PlaceName mismatch.")
	}
}
