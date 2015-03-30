package steps

import (
	"net/http"
	"strconv"

	"github.com/keima/machitools/util"
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

	item := StepItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.Author = u.String()

	if err := item.Save(c); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
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
