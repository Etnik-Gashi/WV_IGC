package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"
	igc "github.com/marni/goigc"
	//"math/rand"
	"encoding/json"
    //"strconv"
)
type metaInfo struct {
	Uptime string `json:"uptime"`
	Info string `json:"info"`
	Version string `json:"version"`

}

var timeStarted = time.Now()

type _url struct {
	URL string `json:"url"`
}


func handler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		temp,err:=template.ParseFiles("wv.html")
		if err !=nil{
			http.Error(w,"Error",http.StatusInternalServerError)
		}
		temp.Execute(w,nil)
		break
	case http.MethodPost:
		input:=r.PostFormValue("inputi")
		pattern:=".*.igc"
		res,err:=regexp.MatchString(pattern,input)
		if err != nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		if res{
			URL := &_url{}
			URL.URL = r.FormValue("url")
			// var jsonR map[string]string
			var _ = json.NewDecoder(r.Body).Decode(URL)

			track,_ := igc.ParseLocation(URL.URL)
			response := "{"
			response += "\"id\": " + "\"" + track.UniqueID + "\","
			response += "\"url\": " + "\"" + URL.URL + "\""
			response += "}"


			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, response)


			return
		}else {
			fmt.Println("Invalid file format, only IGC file!!")
		}
		break
	default:
		http.Error(w,"Not implemented",http.StatusNotImplemented)


	}


}
func handler2(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	uptime := FormatSince(timeStarted)


	metaStruct := &metaInfo{}

	metaStruct.Uptime = uptime
	metaStruct.Info = "Service for igc tracks"
	metaStruct.Version = "v1"

	json.NewEncoder(w).Encode(metaStruct)
}

func main() {

	http.HandleFunc("/igcinfo/api",handler2)
	http.HandleFunc("/igcinfo/api/igc",handler)
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
