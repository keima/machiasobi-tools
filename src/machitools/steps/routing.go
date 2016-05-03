package steps

import (
	"errors"
	"net/http"
	"strconv"
	"machitools/util"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
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
		rest.Error(w, "Invalid param 'id':"+err.Error(), http.StatusBadRequest)
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

func DeleteStep(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	id, _e := strconv.ParseInt(r.PathParam("id"), 10, 64)
	if _e != nil {
		rest.Error(w, _e.Error(), http.StatusBadRequest)
		return
	}

	item := StepItem{Id: id}
	if err := item.Delete(c); err != nil {
		if err == datastore.ErrNoSuchEntity {
			rest.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Errorf(c, "DeleteMenu: %v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
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

	errTx := datastore.RunInTransaction(c, func(tc context.Context) error {
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
		log.Errorf(c, "Tx Error: %v", errTx)
		rest.Error(w, errTx.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
