package machitools

import (
	"github.com/ant0ine/go-json-rest/rest"

	"net/http"

	"log"
	"github.com/keima/machitools/news"
	"github.com/keima/machitools/traffic"
)

func init() {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		PreRoutingMiddlewares: []rest.Middleware{
			&MyMiddleware{},
		},
	}

	trafficItem := traffic.TrafficItem{}
	newsItem := news.NewsItem{}

	err := handler.SetRoutes(
		// Traffic
		rest.RouteObjectMethod("POST", "/api/#version/traffic/#traffic", &trafficItem, "PostTraffic"),
		rest.RouteObjectMethod("GET",  "/api/#version/traffic/#traffic/#direction", &trafficItem, "GetTraffic"),

		// News
		rest.RouteObjectMethod("GET",  "/api/#version/news", &newsItem, "GetNewsList"),
		rest.RouteObjectMethod("GET",  "/api/#version/news/#id", &newsItem, "GetNews"),
		rest.RouteObjectMethod("POST", "/api/#version/news/#id", &newsItem, "PostNews"),

		// Auth
		&rest.Route{"GET", "/api/#version/auth/check", CheckStatus},
		&rest.Route{"GET", "/api/#version/auth/login", Login},
		&rest.Route{"GET", "/api/#version/auth/logout", Logout},
	)

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &handler)
}
