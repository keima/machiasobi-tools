package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

type MyMiddleware struct{}

func (MyMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {

		for key,val := range r.PathParams {
			log.Println(key + " : " + val)
		}

		handler(w, r)
	}
}
