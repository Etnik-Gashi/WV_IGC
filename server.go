package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

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
			fmt.Fprintln(w,"File format is correct")
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
