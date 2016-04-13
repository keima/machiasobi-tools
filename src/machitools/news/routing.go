// Package news はnews(簡易ブログ)を提供するパッケージです
package news

import (
	"net/http"
	"strconv"

	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
	"appengine/user"
)

// GetNewsList では登録されているニュースを取得します
//
// URL Parameter:
// first DB登録されているリストのどの場所から取得開始するか(デフォルトは 0)
//       (datastoreの都合、正確な場所を取れるかどうかは怪しい)
// size  何件取得するか(デフォルトは10)
// private Publicでないものも表示するか（ログインしていないと有効化できない）
//
// Response:
// json: NewsItemのリスト
func GetNewsList(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	first, err := strconv.Atoi(r.FormValue("first"))
	if err != nil {
		first = 0
	}

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		size = 10
	}

	onlyPublic := true
	if r.Request.FormValue("private") == "true" {
		if u := user.Current(c); u != nil && user.IsAdmin(c) {
			onlyPublic = false
		}
	}

	itemList, err := LoadAll(c, first, size, onlyPublic)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&itemList)
}

// PostNews ではニュースを新規登録します。認証必須です。
//
// JSON Parameter:
// Title
// Article
// IsPublic 公開するかどうか
func PostNews(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	news := NewsItem{}

	if err := r.DecodeJsonPayload(&news); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	news.Author = u.String()

	if err := news.Save(c, r.PathParam("id")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&news)
}

// PutNews ではニュースを上書きします。認証必須です。
func PutNews(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	// load old news entity
	oldNews := NewsItem{}
	err := oldNews.Load(c, r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set news json to news entity
	newNews := NewsItem{}
	if err := r.DecodeJsonPayload(&newNews); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// merge new news to old news
	oldNews.Title = newNews.Title
	oldNews.Article = newNews.Article
	oldNews.IsPublic = newNews.IsPublic
	oldNews.Author = u.String()

	if err := oldNews.Save(c, r.PathParam("id")); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&oldNews)
}

// GetNews では登録されているニュースから指定されたIDのものだけ返却します
//
// Path Parameter:
// id: ニュースのID
func GetNews(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	news := NewsItem{}

	err := news.Load(c, r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&news)
}
