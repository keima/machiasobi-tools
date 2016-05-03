package machitools

import (
	"log"
	"net/http"
	"regexp"
	"google.golang.org/appengine"
	"gopkg.in/ant0ine/go-json-rest.v2/rest"

	"machitools/delay"
	"machitools/event"
	"machitools/maps"
	"machitools/news"
	"machitools/steps"
	"machitools/traffic"
	"machitools/weather"
	"machitools/periods"
	"machitools/favorite"
	"machitools/calendar"
	"machitools/menu"
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
						if m, _ := regexp.MatchString("^/api/\\w*/(auth|favorite|periods|calendars)", request.URL.Path); !m {
							// マチアプリAPIはGETリクエストでもoriginチェックさせたいのでtrue返さない
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
		&rest.Route{"PUT", pathPrefix + "/news/:id", news.PutNews},

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
		&rest.Route{"DELETE", pathPrefix + "/steps/:id", steps.DeleteStep},

		// Weather
		&rest.Route{"GET", pathPrefix + "/weather/:id", weather.GetWeather},

		// Calendar
		&rest.Route{"GET", pathPrefix + "/calendars", calendar.GetCalendarList},
		&rest.Route{"POST", pathPrefix + "/calendars", calendar.PostCalendar},
		&rest.Route{"POST", pathPrefix + "/calendars/order", calendar.PostOrder},
		&rest.Route{"GET", pathPrefix + "/calendars/:id", calendar.GetCalendar},
		&rest.Route{"PUT", pathPrefix + "/calendars/:id", calendar.PutCalendar},
		&rest.Route{"DELETE", pathPrefix + "/calendars/:id", calendar.DeleteCalendar},

		// Favorite
		&rest.Route{"GET", pathPrefix + "/favorite", favorite.GetFavList},
		&rest.Route{"POST", pathPrefix + "/favorite/#calId/#eventId", favorite.PostFav},
		&rest.Route{"DELETE", pathPrefix + "/favorite/#calId/#eventId", favorite.DeleteFav},

		// Periods
		&rest.Route{"GET", pathPrefix + "/periods", periods.GetPeriodList},
		&rest.Route{"POST", pathPrefix + "/periods", periods.PostPeriod},
		&rest.Route{"POST", pathPrefix + "/periods/:id/deactivate", periods.DeActivatePeriod},

		// Auth
		&rest.Route{"GET", pathPrefix + "/auth/check", CheckStatus},
		&rest.Route{"GET", pathPrefix + "/auth/login", Login},
		&rest.Route{"GET", pathPrefix + "/auth/logout", Logout},

		// Menu
		&rest.Route{"GET", pathPrefix + "/menu", menu.GetMenuList},
		&rest.Route{"POST", pathPrefix + "/menu", menu.PostMenu},
		&rest.Route{"POST", pathPrefix + "/menu/order", menu.PostOrder},
		&rest.Route{"GET", pathPrefix + "/menu/:id", menu.GetMenu},
		&rest.Route{"PUT", pathPrefix + "/menu/:id", menu.PutMenu},
		&rest.Route{"DELETE", pathPrefix + "/menu/:id", menu.DeleteMenu},

	)
	//@formatter:on

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
