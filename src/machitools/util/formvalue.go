package util

import (
	"net/http"
	"strconv"
	"google.golang.org/appengine/user"
	"golang.org/x/net/context"
)

func ParseFirstSizePrivate(c context.Context, r *http.Request) (int, int, bool) {
	first, err := strconv.Atoi(r.FormValue("first"))
	if err != nil {
		first = 0
	}

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		size = 10
	}

	private := false
	if r.FormValue("private") == "true" {
		if u := user.Current(c); u != nil && user.IsAdmin(c) {
			private = true
		}
	}

	return first, size, private
}
