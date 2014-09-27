package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

type MyMiddleware struct{}

func (MyMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {

		corsInfo := request.GetCorsInfo()

		// Be nice with non CORS requests, continue
		// Alternatively, you may also chose to only allow CORS requests, and return an error.
		if !corsInfo.IsCors {
			// continue, execute the wrapped middleware
			handler(writer, request)
			return
		}

		if request.Method == "GET" {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			handler(writer, request)
			return
		} else {
			// Probably, error string will hidden by browser.
			rest.Error(writer, "Invalid Origin", http.StatusForbidden)
			return
		}

	}
}
