// Package news はnews(簡易ブログ)を提供するパッケージです
package news

import (
	"appengine"
	"appengine/user"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strconv"
)

// 登録されているニュースを指定件数取得します
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

// ニュースを登録します。認証必須です。
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
