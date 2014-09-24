package traffic

import (
	"github.com/ant0ine/go-json-rest/rest"
	"appengine"
	"net/http"
	"appengine/user"
	"time"
)

func (item *TrafficItem) GetTraffic(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	trafficName := TrafficName(r.PathParam("traffic"))
	if trafficName == "" {
		rest.Error(w, "Traffic is allowed bus or ropeway.", http.StatusBadRequest)
		return
	}

	direction := Direction(r.PathParam("direction"))
	if direction == DirectionError {
		rest.Error(w, "Direction is allowed inbound or outbound.", http.StatusBadRequest)
		return
	}

	i, err := item.loadLatest(c, trafficName, direction)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(i)
}

func (item *TrafficItem) PostTraffic(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	// check pathParam
	if TrafficName(r.PathParam("traffic")) == "" || Direction(r.PathParam("direction")) == DirectionError {
		rest.Error(w, "PathParam Error.", http.StatusBadRequest)
		return
	}

	traffic := TrafficItem{}

	if err := r.DecodeJsonPayload(&traffic); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	traffic.Direction = Direction(r.PathParam("direction"))
	traffic.Author = u.String()
	traffic.Date = time.Now()

	if _,err := traffic.save(c, TrafficName(r.PathParam("traffic"))); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&traffic)
}

func TrafficName(trafficName string) string{
	switch trafficName {
	case "bus", "ropeway":
		return trafficName
	default:
		return ""
	}
}

func Direction(directionName string) int {
	switch directionName {
	case "inbound":
		return DirectionInbound
	case "outbound":
		return DirectionOutbound
	default:
		return DirectionError
	}
}
