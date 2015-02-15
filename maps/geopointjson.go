package maps

import (
	"encoding/json"

	"appengine"
)

type GeoPointJson appengine.GeoPoint

func (g *GeoPointJson) UnmarshalJSON(b []byte) (err error) {
	var jm map[string]float64
	if err = json.Unmarshal(b, &jm); err == nil {
		g.Lat = jm["latitude"]
		g.Lng = jm["longitude"]
	}
	return
}

func (g *GeoPointJson) MarshalJSON() ([]byte, error) {
	jm := make(map[string]float64)
	jm["latitude"] = g.Lat
	jm["longitude"] = g.Lng
	return json.Marshal(jm)
}
