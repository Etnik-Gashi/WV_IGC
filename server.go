package main

import (
	"fmt"
	//"html/template"
	"net/http"
	"regexp"
	"time"
	igc "github.com/marni/goigc"
	"math/rand"
	"encoding/json"
    "strconv"
)


var timeStarted = time.Now()

type _url struct {
	URL string `json:"url"`
}
var igcFiles []Track


type Track struct {
	Id string   `json:"id"`
	igcTrack igc.Track `json:"igc_track"`
}
func handler(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "{" + "\"uptime\": \""+FormatSince(timeStarted)+"\"," + "\"info\": \"Service for IGC tracks.\"," + "\"version\": \"v1\""+ "}")
}

func handler2(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		ids := make([]string, 0, 0)

		for i := range igcFiles {
			ids = append(ids, igcFiles[i].Id)
		}

		json.NewEncoder(w).Encode(ids)

		break
	case http.MethodPost:

		pattern:=".*.igc"

		w.Header().Set("Content-Type", "application/json")
		//jsonR := make(map[string]string)
		URL := &_url{}

		var error = json.NewDecoder(r.Body).Decode(URL)
		if error != nil {
			fmt.Fprintln(w, "Error!! ", error)
			return
		}
		res,err:=regexp.MatchString(pattern,URL.URL)
		if err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		if res {

			track, _ := igc.ParseLocation(URL.URL)

			Id := rand.Intn(1000)

			igcFile := Track{}
			igcFile.Id = strconv.Itoa(Id)
			igcFile.igcTrack = track

			igcFiles = append(igcFiles, igcFile)

			json.NewEncoder(w).Encode(igcFile.Id)
			return
		}
		break
	default:
		http.Error(w,"Not implemented",http.StatusNotImplemented)


	}


}





func main() {

	http.HandleFunc("/igcinfo/api/",handler)
	http.HandleFunc("/igcinfo/api/igc/",handler2)
	http.ListenAndServe(":8080",nil)
}
func FormatSince(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	f := ts / Decisecond
	y := d / 365
	return fmt.Sprintf("P%dY%dD%dH%dM%d.%dS", y, d, h, m, s, f)
}
