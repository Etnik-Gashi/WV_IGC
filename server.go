package main

import (
	"fmt"
	"log"
	"os"
	//"html/template"
	"net/http"
	"regexp"
	"time"
	igc "github.com/marni/goigc"
	"math/rand"
	"encoding/json"
	"strconv"
	//"path/filepath"
	"strings"
)


var timeStarted = time.Now()

//IgcFiles is a slice for storing igc files
var igcFiles []Track


type _url struct {
	URL string `json:"url"`
}
//Calculating the total length of track
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

//Track structure: This is structure for storing the track and to access their's id
type Track struct {
	ID string   `json:"id"`
	IgcTrack igc.Track `json:"igc_track"`
}
//Attributes : the info about each igc file via id
type Attributes struct{
	HeaderDate string `json:"h_date"`
	Pilot string `json:"pilot"`
	Glider string `json:"glider"`
	GliderID string 	`json:"glider_id"`
	Length float64 `json:"track_length"`
}
//Calculating uptime based on ISO 8601
func timeSince(t time.Time) string {

	Decisecond := 100 * time.Millisecond
	Day        := 24 * time.Hour

	ts := time.Since(t)
	sign := time.Duration(1)

	ts += Decisecond / 2
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


//Handling based in parsing url
func handler(w http.ResponseWriter,r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	//var empty = regexp.MustCompile(``)
	var api= regexp.MustCompile(`api`)
	switch {
	//Handling for /igcinfo and for /<rubbish>
	case len(parts) == 2 :
		http.Error(w, "404 - Page not found!", http.StatusNotFound)
		break

		//Handling for /igcinfo/api
	case len(parts) == 3 && api.MatchString(parts[2]) :
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{"+"\"uptime\": \""+timeSince(timeStarted)+"\","+"\"info\": \"Service for IGC tracks.\","+"\"version\": \"v1\""+"}")
		break
		//Handling for /igcinfo/api/igc
	case len(parts) == 4 :
		{

			w.Header().Set("Content-Type", "application/json")
			var rNum= regexp.MustCompile(`/igcinfo/api/igc`)

			switch {

			case rNum.MatchString(r.URL.Path):

				switch r.Method {
				//Handling GET /igcinfo/api/igc for returning all ids storing in a slice
				case http.MethodGet:
					ids := make([]string, 0, 0)

					for i := range igcFiles {
						ids = append(ids, igcFiles[i].ID)
					}

					json.NewEncoder(w).Encode(ids)

					break
				case http.MethodPost:

					//handling post /igcinfo/api/igc for sending a url and returning an id for that url
					pattern := ".*.igc"

					URL := &_url{}

					var error= json.NewDecoder(r.Body).Decode(URL)
					if error != nil {
						fmt.Fprintln(w, "Error!! ", error)
						return
					}
					res, err := regexp.MatchString(pattern, URL.URL)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						fmt.Fprintln(w, "Error!! ", error)
						return
					}
					if res {

						track, _ := igc.ParseLocation(URL.URL)

						ID := rand.Intn(1000)

						igcFile := Track{}
						igcFile.ID = strconv.Itoa(ID)
						igcFile.IgcTrack = track

						igcFiles = append(igcFiles, igcFile)

						json.NewEncoder(w).Encode(igcFile.ID)
						return
					}
					break
				default:
					http.Error(w, "Not implemented", http.StatusNotImplemented)

				}
			default:
				http.Error(w, "404 - Page not found!", http.StatusNotFound)
				break
			}

		}
		break
	case len(parts) == 5 : {
		//Handling /igcinfo/api/igc/<id>

		w.Header().Set("Content-Type", "application/json")

		attributes := &Attributes{}

		var rNum= regexp.MustCompile(`/igcinfo/api/igc/\d{1,}`)

		switch {
		case rNum.MatchString(r.URL.Path):

			for i := range igcFiles {

				if igcFiles[i].ID == parts[4] {
					attributes.HeaderDate = igcFiles[i].IgcTrack.Header.Date.String()
					attributes.Pilot = igcFiles[i].IgcTrack.Pilot
					attributes.Glider = igcFiles[i].IgcTrack.GliderType
					attributes.GliderID = igcFiles[i].IgcTrack.GliderID
					attributes.Length = trackLength(igcFiles[i].IgcTrack)

					json.NewEncoder(w).Encode(attributes)
				}
				//Handling if user type different id from ids stored
				if igcFiles[i].ID != parts[4]{
					http.Error(w,"", http.StatusNotFound)
				}

			}

			break
		default:
			http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)

		}
	}
		break

		//Handling for GET /api/igc/<id>/<field>
	case len(parts) == 6 :{
		var rNum= regexp.MustCompile(`/igcinfo/api/igc/\d{1,}/\w{1,}`)

		switch {

		case rNum.MatchString(r.URL.Path):

			for i := range igcFiles {

				if igcFiles[i].ID == parts[4] {
					switch {
					case parts[5] == "pilot":
						fmt.Fprintln(w,igcFiles[i].IgcTrack.Pilot)
						break
					case parts[5] == "glider":
						fmt.Fprintln(w,igcFiles[i].IgcTrack.GliderType)
						break
					case parts[5] == "glider_id":
						fmt.Fprintln(w,igcFiles[i].IgcTrack.GliderID)
						break
					case parts[5] == "track_length":
						fmt.Fprintln(w,trackLength(igcFiles[i].IgcTrack))
						break
					case parts[5] == "h_date":
						fmt.Fprintln(w,igcFiles[i].IgcTrack.Header.Date.String())
						break
					default:
						http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
						break
					}

				}

			}
			break
		default:
			http.Error(w, "400 - Bad Request", http.StatusBadRequest)

		}

	}
		break
	case len(parts)>6:{
		http.Error(w, "404 - Page not found!", http.StatusNotFound)

	}
		break

	default:
		http.Error(w, "404 - Page not found!", http.StatusNotFound)
		break
	}
}

func main() {

	http.HandleFunc("/",handler)

	//fmt.Println("listening...")

	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

