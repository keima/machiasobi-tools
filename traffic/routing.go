package traffic

import (
	"github.com/ant0ine/go-json-rest/rest"
	"appengine"
	"net/http"
	"appengine/user"
)

func GetTraffic(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	trafficName := trafficName(r.PathParam("traffic"))
	if trafficName == "" {
		rest.Error(w, "Traffic is allowed bus or ropeway.", http.StatusBadRequest)
		return
	}

	direction := direction(r.PathParam("direction"))
	if direction == DirectionError {
		rest.Error(w, "Direction is allowed inbound or outbound.", http.StatusBadRequest)
		return
	}

	item, err := LoadLatest(c, trafficName, direction)
	if err != nil {
		if err == ErrItemNotFound {
			rest.Error(w, err.Error(), http.StatusNotFound)
		} else {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteJson(item)
}

func PostTraffic(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)

	u := user.Current(c)
	if u == nil || !user.IsAdmin(c) {
		rest.Error(w, "Administrator login Required.", http.StatusUnauthorized)
		return
	}

	// check pathParam
	if trafficName(r.PathParam("traffic")) == "" || direction(r.PathParam("direction")) == DirectionError {
		rest.Error(w, "PathParam Error.", http.StatusBadRequest)
		return
	}

	traffic := TrafficItem{}

	if err := r.DecodeJsonPayload(&traffic); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	traffic.Direction = direction(r.PathParam("direction"))
	traffic.Author = u.String()

	if err := traffic.Save(c, trafficName(r.PathParam("traffic"))); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&traffic)
}

func trafficName(name string) string{
	switch name {
	case "bus", "ropeway":
		return name
	default:
		return ""
	}
}

func direction(directionName string) int {
	switch directionName {
	case "inbound":
		return DirectionInbound
	case "outbound":
		return DirectionOutbound
	default:
		return DirectionError
	}
}
