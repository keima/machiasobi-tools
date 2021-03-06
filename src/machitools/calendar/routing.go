package calendar

import (
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"github.com/mjibson/goon"
)

func GetCalendarList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	builder := NewCalendarItemQueryBuilder()
	builder.Order.Asc()

	visibility := r.FormValue("visibility")
	if !user.IsAdmin(c) || visibility != "all" {
		builder.Enabled.Equal(true)
	}

	items := CalendarItemList{}

	if _,err := goon.FromContext(c).GetAll(builder.Query(), &items); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if visibility != "all" {
		w.Header().Set("cache-control", "private, max-age=900") // 15min
	}

	w.WriteJson(&items)
}

// 管理者専用API
func PostCalendar(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := CalendarItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var e error
	item.Order, e = CalendarItemList{}.Count(r.Request)
	if e != nil {
		rest.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	if err := item.Save(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// 管理者専用API
func PostOrder(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	ids := []string{}
	if err := r.DecodeJsonPayload(&ids); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errTx := datastore.RunInTransaction(c, func(tc context.Context) error {
		for i, id := range ids {
			item := CalendarItem{}

			if err := item.Load(r.Request, id); err != nil {
				return err
			}

			item.Order = i

			if err := item.Save(r.Request); err != nil {
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

func GetCalendar(w rest.ResponseWriter, r *rest.Request) {
	idStr := r.PathParam("id")
	item := CalendarItem{}

	if err := item.Load(r.Request, idStr); err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteJson(&item)
}

// 管理者専用API
func PutCalendar(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	idStr := r.PathParam("id")
	existItem := CalendarItem{}
	if err := existItem.Load(r.Request, idStr); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := CalendarItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := item.Save(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// 管理者専用API
func DeleteCalendar(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	idStr := r.PathParam("id")
	existItem := CalendarItem{}
	if err := existItem.Load(r.Request, idStr); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := existItem.Delete(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}