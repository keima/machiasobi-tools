package machitools

import (
	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"net/http"

	"log"

	"github.com/keima/machitools/delay"
	"github.com/keima/machitools/event"
	"github.com/keima/machitools/maps"
	"github.com/keima/machitools/news"
	"github.com/keima/machitools/steps"
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

	//@formatter:off
	err := handler.SetRoutes(
		// Traffic
		&rest.Route{"POST", PathPrefix + "/traffic/:traffic/:direction", traffic.PostTraffic},
		&rest.Route{"GET", PathPrefix + "/traffic/:traffic/:direction", traffic.GetTraffic},

		// Ticket
		&rest.Route{"GET", PathPrefix + "/events", event.GetEventList},
		&rest.Route{"POST", PathPrefix + "/events", event.PostEvent},
		&rest.Route{"GET", PathPrefix + "/events/:id", event.GetEvent},
		&rest.Route{"POST", PathPrefix + "/events/:id", event.PostEvent},
		//		&rest.Route{"POST", PathPrefix + "/events/:id/done",   event.GetEvent},
		//		&rest.Route{"POST", PathPrefix + "/events/:id/delete", event.GetEvent},

		// Delay
		&rest.Route{"GET", PathPrefix + "/delay/:place", delay.GetDelay},
		&rest.Route{"POST", PathPrefix + "/delay/:place", delay.PostDelay},

		// News
		&rest.Route{"GET", PathPrefix + "/news", news.GetNewsList},
		&rest.Route{"GET", PathPrefix + "/news/:id", news.GetNews},
		&rest.Route{"POST", PathPrefix + "/news/:id", news.PostNews},

		// Maps
		&rest.Route{"GET", PathPrefix + "/maps", maps.GetMapList},
		&rest.Route{"GET", PathPrefix + "/maps/:id", maps.GetMap},
		&rest.Route{"POST", PathPrefix + "/maps/:id", maps.PostMap},
		&rest.Route{"PUT", PathPrefix + "/maps/:id", maps.PutMap},
		&rest.Route{"POST", PathPrefix + "/maps/:id/markers", maps.PostMarker},
		&rest.Route{"DELETE", PathPrefix + "/maps/:id/markers/:key", maps.DeleteMarker},

		// Steps
		&rest.Route{"GET", PathPrefix + "/steps", steps.GetStepList},
		&rest.Route{"POST", PathPrefix + "/steps", steps.PostStep},
		&rest.Route{"POST", PathPrefix + "/steps/order", steps.PostOrder},
		&rest.Route{"GET", PathPrefix + "/steps/:id", steps.GetStep},
		&rest.Route{"POST", PathPrefix + "/steps/:id", steps.UpdateStep},

		// Auth
		&rest.Route{"GET", PathPrefix + "/auth/check", CheckStatus},
		&rest.Route{"GET", PathPrefix + "/auth/login", Login},
		&rest.Route{"GET", PathPrefix + "/auth/logout", Logout},
	)
	//@formatter:on

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
