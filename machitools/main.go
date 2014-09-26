package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"

	"net/http"

	"log"
	"github.com/keima/machitools/news"
	"github.com/keima/machitools/traffic"
)

const PathPrefix = "/api/#version"

func init() {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		PreRoutingMiddlewares: []rest.Middleware{
			&MyMiddleware{},
		},
	}

	err := handler.SetRoutes(
		// Traffic
		&rest.Route{"POST", PathPrefix + "/traffic/:traffic/:direction", traffic.PostTraffic},
		&rest.Route{"GET",  PathPrefix + "/traffic/:traffic/:direction", traffic.GetTraffic},

		// Ticket
		/*
		&rest.Route{"GET",  PathPrefix + "/tickets/list",       &trafficItem, "GetTraffic"},
		&rest.Route{"GET",  PathPrefix + "/tickets/:id",        &trafficItem, "GetTraffic"},
		&rest.Route{"POST", PathPrefix + "/tickets/:id/update", &trafficItem, "GetTraffic"},
		&rest.Route{"POST", PathPrefix + "/tickets/:id/done",   &trafficItem, "GetTraffic"},
		&rest.Route{"POST", PathPrefix + "/tickets/:id/delete", &trafficItem, "GetTraffic"},
		*/

		// News
		&rest.Route{"GET",  PathPrefix + "/news",     news.GetNewsList},
		&rest.Route{"GET",  PathPrefix + "/news/:id", news.GetNews},
		&rest.Route{"POST", PathPrefix + "/news/:id", news.PostNews},

		// Auth
		&rest.Route{"GET", PathPrefix + "/auth/check",  CheckStatus},
		&rest.Route{"GET", PathPrefix + "/auth/login",  Login},
		&rest.Route{"GET", PathPrefix + "/auth/logout", Logout},
	)

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
