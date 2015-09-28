package maps

import (
	"net/http"
	"strconv"

	"github.com/keima/machiasobi-tools/util"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Mapの一覧を取得します。公開APIです。
func GetMapList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	first, size, private := util.ParseFirstSizePrivate(c, r.Request)

	itemList, err := LoadAll(c, first, size, private)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&itemList)
}

// Mapを取得します。公開APIです。
func GetMap(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	keyName := r.PathParam("id")
	if keyName == "" {
		rest.Error(w, "Invalid param 'id'", http.StatusInternalServerError)
		return
	}

	item := Map{}

	if err := item.Load(c, keyName); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)
	if items, err := LoadAllMapItem(c, key); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(*items) > 0 {
		item.MapItems = *items
	}

	w.WriteJson(&item)
}

// Mapを登録します。管理者権限必須です。
func PostMap(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := Map{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := item.Save(c, r.PathParam("id")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// Mapを更新します。管理者権限必須です。
func PutMap(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	item := Map{}
	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := item.Save(c, r.PathParam("id")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

// MapMarkerを登録します。管理者権限必須です。
func PostMarker(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	keyName := r.PathParam("id")
	if keyName == "" {
		rest.Error(w, "Invalid param 'id'", http.StatusInternalServerError)
		return
	}

	item := MapItem{}

	if err := r.DecodeJsonPayload(&item); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)
	if err := item.Save(c, key); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&item)
}

func DeleteMarker(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	keyName := r.PathParam("id")
	if keyName == "" {
		rest.Error(w, "Invalid param 'id'", http.StatusInternalServerError)
		return
	}

	mapItemKeyName := r.PathParam("key")
	if mapItemKeyName == "" {
		rest.Error(w, "Invalid param 'key'", http.StatusInternalServerError)
		return
	}

	keyId, err := strconv.ParseInt(mapItemKeyName, 10, 64)
	if err != nil {
		rest.Error(w, "Invalid param 'key'", http.StatusInternalServerError)
		return
	}

	parent := datastore.NewKey(c, kindNameMap, keyName, 0, nil)
	key := datastore.NewKey(c, kindNameMapItem, "", keyId, parent)
	if err := DeleteMapItem(c, key); err != nil {
		rest.Error(w, "id or key is not match (or DB error).", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
