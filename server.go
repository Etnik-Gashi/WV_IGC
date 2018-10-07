package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var timeStarted = time.Now()

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
			v := FormatSince(timeStarted)
			fmt.Fprintf(w,"File format is correct \n Time: %s",v)
			return
		}else {
			fmt.Println("Invalid file format, only IGC file!!")
		}
		break
	default:
		http.Error(w,"Not implemented",http.StatusNotImplemented)


	}


}

func main() {
	http.HandleFunc("/igcinfo",handler)
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
