package machitools

import (
	"log"
	"net/http"
	"regexp"

	"github.com/keima/machitools/calendar"
	"github.com/keima/machitools/delay"
	"github.com/keima/machitools/event"
	"github.com/keima/machitools/maps"
	"github.com/keima/machitools/news"
	"github.com/keima/machitools/steps"
	"github.com/keima/machitools/traffic"
	"github.com/keima/machitools/weather"

	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"appengine"
)

const pathPrefix = "/api/#version"

func init() {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.CorsMiddleware{
				RejectNonCorsRequests: false,
				OriginValidator: func(origin string, request *rest.Request) bool {
					if appengine.IsDevAppServer() {
						return true
					} else if request.Method == "GET" {
						if m, _ := regexp.MatchString("^/auth", str); !m {
							// 認証系はGETリクエストであれoriginチェックさせたいのでtrue返さない
							return true
						}
					}
					return origin == "http://machi.p-side.net"
				},
				AllowedMethods:                []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:                []string{"Accept", "Content-Type", "X-Custom-Header", "Origin"},
				AccessControlAllowCredentials: true,
				AccessControlMaxAge:           3600,
			},
		},
	}

	//@formatter:off
	err := handler.SetRoutes(
		// Traffic
		&rest.Route{"POST", pathPrefix + "/traffic/:traffic/:direction", traffic.PostTraffic},
		&rest.Route{"GET", pathPrefix + "/traffic/:traffic/:direction", traffic.GetTraffic},

		// Ticket
		&rest.Route{"GET", pathPrefix + "/events", event.GetEventList},
		&rest.Route{"POST", pathPrefix + "/events", event.PostEvent},
		&rest.Route{"GET", pathPrefix + "/events/:id", event.GetEvent},
		&rest.Route{"POST", pathPrefix + "/events/:id", event.PostEvent},
		//		&rest.Route{"POST", pathPrefix + "/events/:id/done",   event.GetEvent},
		//		&rest.Route{"POST", pathPrefix + "/events/:id/delete", event.GetEvent},

		// Delay
		&rest.Route{"GET", pathPrefix + "/delay/:place", delay.GetDelay},
		&rest.Route{"POST", pathPrefix + "/delay/:place", delay.PostDelay},

		// News
		&rest.Route{"GET", pathPrefix + "/news", news.GetNewsList},
		&rest.Route{"GET", pathPrefix + "/news/:id", news.GetNews},
		&rest.Route{"POST", pathPrefix + "/news/:id", news.PostNews},

		// Maps
		&rest.Route{"GET", pathPrefix + "/maps", maps.GetMapList},
		&rest.Route{"GET", pathPrefix + "/maps/:id", maps.GetMap},
		&rest.Route{"POST", pathPrefix + "/maps/:id", maps.PostMap},
		&rest.Route{"PUT", pathPrefix + "/maps/:id", maps.PutMap},
		&rest.Route{"POST", pathPrefix + "/maps/:id/markers", maps.PostMarker},
		&rest.Route{"DELETE", pathPrefix + "/maps/:id/markers/:key", maps.DeleteMarker},

		// Steps
		&rest.Route{"GET", pathPrefix + "/steps", steps.GetStepList},
		&rest.Route{"POST", pathPrefix + "/steps", steps.PostStep},
		&rest.Route{"POST", pathPrefix + "/steps/order", steps.PostOrder},
		&rest.Route{"GET", pathPrefix + "/steps/:id", steps.GetStep},
		&rest.Route{"POST", pathPrefix + "/steps/:id", steps.UpdateStep},

		// Weather
		&rest.Route{"GET", pathPrefix + "/weather/:id", weather.GetWeather},

		// Calendar
		&rest.Route{"GET", pathPrefix + "/calendar", calendar.GetFavList},
		&rest.Route{"POST", pathPrefix + "/calendar/#calId/#eventId", calendar.PostFav},
		&rest.Route{"DELETE", pathPrefix + "/calendar/#calId/#eventId", calendar.DeleteFav},

		// Auth
		&rest.Route{"GET", pathPrefix + "/auth/check", CheckStatus},
		&rest.Route{"GET", pathPrefix + "/auth/login", Login},
		&rest.Route{"GET", pathPrefix + "/auth/logout", Logout},
	)
	//@formatter:on

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
