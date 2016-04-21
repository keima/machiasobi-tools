package menu

import (
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"github.com/mjibson/goon"
	"net/http"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// GetMenuList is Public API
func GetMenuList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	builder := NewMenuItemQueryBuilder()
	builder.OrderIndex.Asc()

	vis := r.FormValue("visibility")
	if user.IsAdmin(c) && vis == "all" {
		builder.Enabled.Equal(true)
	}

	items := MenuList{}
	if _, err := goon.FromContext(c).GetAll(builder.Query(), &items); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("cache-control", "private, max-age=900") // 15min
	w.WriteJson(&items)
}

func PostMenu(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := MenuItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	g := goon.FromContext(c);
	builder := NewMenuItemQueryBuilder()

	// orderをアイテム群の末尾となるようにする
	if order, err := g.Count(builder.Query()); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		item.OrderIndex = order;
	}

	if _, err := g.Put(item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&item)
}

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

	g := goon.FromContext(c)
	if err := g.RunInTransaction(func(tg *goon.Goon) error {
		for i, id := range ids {
			item := MenuItem{Id: id}
			if err := tg.Get(item); err != nil {
				return err
			}

			item.OrderIndex = i
			if _, err := tg.Put(item); err != nil {
				return err
			}
		}
		return nil
	}, &datastore.TransactionOptions{XG: true}); err != nil {
		log.Errorf(c, "Tx Error: %v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	GetMenuList(w, r)
}

func GetMenu(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	item := MenuItem{Id: r.PathParam("id")}
	if err := goon.FromContext(c).Get(item); err != nil {
		if err == datastore.ErrNoSuchEntity {
			rest.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Errorf(c, "GetMenu: %v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteJson(&item)
}

func PutMenu(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	g := goon.FromContext(c);

	item := MenuItem{Id: r.PathParam("id")}
	if err := g.Get(item); err != nil {
		if err == datastore.ErrNoSuchEntity {
			rest.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Errorf(c, "GetMenu: %v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	itemNew := MenuItem{}
	if err := r.DecodeJsonPayload(&itemNew); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemNew.Id = item.Id

	if _, err := g.Put(itemNew); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&itemNew)
}

func DeleteMenu(w rest.ResponseWriter, r *rest.Request) {
    c := appengine.NewContext(r.Request)

	if !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	g := goon.FromContext(c);

	item := MenuItem{Id: r.PathParam("id")}

	if err := g.Delete(g.Key(item)); err != nil {
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
