package steps

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/keima/machiasobi-tools/util"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func GetStepList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	first, size, private := util.ParseFirstSizePrivate(c, r.Request)

	items, err := LoadAll(c, first, size, private)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&items)
}

func GetStep(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	keyId, err := strconv.ParseInt(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, "Invalid param 'id':"+err.Error(), http.StatusInternalServerError)
		return
	}

	item := StepItem{}
	if err := item.Load(c, keyId); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func PostStep(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item, errDec := decodeAndValidate(r, u)
	if errDec != nil {
		rest.Error(w, errDec.Error(), http.StatusInternalServerError)
		return
	}

	if err := item.Save(c); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func UpdateStep(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item, errDec := decodeAndValidate(r, u)
	if errDec != nil {
		rest.Error(w, errDec.Error(), http.StatusInternalServerError)
		return
	}

	keyId, err := strconv.ParseInt(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, errDec.Error(), http.StatusInternalServerError)
		return
	}
	if err := item.Update(c, keyId); err != nil {
		rest.Error(w, errDec.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func decodeAndValidate(r *rest.Request, u *user.User) (*StepItem, error) {
	item := StepItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		return nil, err
	}

	matched := false
	for _, v := range AllowedType {
		if item.Type == v {
			matched = true
			break
		}
	}
	if !matched {
		return nil, errors.New("Type `" + item.Type + "` is not allowed.")
	}

	item.Author = u.String()

	return &item, nil
}

func PostOrder(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	ids := []int64{}
	if err := r.DecodeJsonPayload(&ids); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errTx := datastore.RunInTransaction(c, func(c appengine.Context) error {
		for i, id := range ids {
			item := StepItem{}

			if err := item.Load(c, id); err != nil {
				return err
			}

			item.Order = i

			if err := item.Update(c, id); err != nil {
				return err
			}
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})

	if errTx != nil {
		c.Errorf("Tx Error: %v", errTx)
		rest.Error(w, errTx.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
