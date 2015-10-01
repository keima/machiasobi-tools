package periods
import (
	"net/http"
	"appengine"
	"appengine/user"
	"strconv"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	"time"
)

// 公開API
func GetPeriodList(w rest.ResponseWriter, r *rest.Request) {
	items := PeriodItemList{}

	err := items.LoadActive(r.Request)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r.Request)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	for _,item := range items {
		item.Date = item.Date.In(loc)

		// 管理者ユーザーでないときはIDを消す
		if !user.IsAdmin(c) {
			item.Id = 0
		}
	}

	w.Header().Set("cache-control", "private, max-age=3600")
	w.WriteJson(&items)
}

// 管理者専用API
func PostPeriod(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)

	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := PeriodItem{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	item.Date = item.Date.In(loc)
	item.IsActive = true

	if err := item.Save(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// 管理者専用API
func DeActivatePeriod(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)

	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	idStr := r.PathParam("id")
	if idStr == "" {
		rest.Error(w, "param `id` is empty.", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64);
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item := PeriodItem{Id:id}
	if err := item.Load(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.IsActive = false
	if err := item.Save(r.Request); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}