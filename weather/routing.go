package weather

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/ant0ine/go-json-rest.v2/rest"
	valid "gopkg.in/asaskevich/govalidator.v1"

	"appengine"
	"appengine/urlfetch"
)

func GetWeather(w rest.ResponseWriter, r *rest.Request) {
	c := appengine.NewContext(r.Request)
	cityId := r.PathParam("id")

	if !valid.IsNumeric(cityId) {
		rest.Error(w, "cityId: "+cityId+" is not allowed.", http.StatusBadRequest)
		return
	}

	apiUrl := "http://weather.livedoor.com/forecast/webservice/json/v1?city=" + cityId

	client := urlfetch.Client(c)
	resp, err := client.Get(apiUrl)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.(http.ResponseWriter).Write(data)
}
