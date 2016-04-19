package menu

import (
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"github.com/mjibson/goon"
	"net/http"
)

// GetMenuList is Public API
func GetMenuList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	builder := NewMenuItemQueryBuilder()

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