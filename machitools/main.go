package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"

	"net/http"

	"github.com/keima/machitools/event"
	"github.com/keima/machitools/news"
	"github.com/keima/machitools/traffic"
	"log"
)

const PathPrefix = "/api/#version"

func init() {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		PreRoutingMiddlewares: []rest.Middleware{
			&MyMiddleware{},
		},
	}

	//@formatter:off
	err := handler.SetRoutes(
		// Traffic
		&rest.Route{"POST", PathPrefix + "/traffic/:traffic/:direction", traffic.PostTraffic},
		&rest.Route{"GET",  PathPrefix + "/traffic/:traffic/:direction", traffic.GetTraffic},

		// Ticket
		&rest.Route{"GET",  PathPrefix + "/events",      event.GetEventList},
		&rest.Route{"POST", PathPrefix + "/events",      event.PostEvent},
		&rest.Route{"GET",  PathPrefix + "/events/:id",  event.GetEvent},
		&rest.Route{"POST", PathPrefix + "/events/:id",  event.PostEvent},
		//		&rest.Route{"POST", PathPrefix + "/events/:id/done",   event.GetEvent},
		//		&rest.Route{"POST", PathPrefix + "/events/:id/delete", event.GetEvent},

		// News
		&rest.Route{"GET",  PathPrefix + "/news",        news.GetNewsList},
		&rest.Route{"GET",  PathPrefix + "/news/:id",    news.GetNews},
		&rest.Route{"POST", PathPrefix + "/news/:id",    news.PostNews},

		// Auth
		&rest.Route{"GET",  PathPrefix + "/auth/check",  CheckStatus},
		&rest.Route{"GET",  PathPrefix + "/auth/login",  Login},
		&rest.Route{"GET",  PathPrefix + "/auth/logout", Logout},
	)
	//@formatter:on

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
