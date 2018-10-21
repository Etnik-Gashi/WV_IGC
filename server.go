package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http" //"html/template"
	"regexp"
	"strconv"
	"strings"
	"time" //"path/filepath"

	"github.com/gorilla/mux"
	igc "github.com/marni/goigc"
)

var timeStarted = time.Now()

var urlMap = make(map[int]string)
var mapID int
var initialID int
var uniqueId int

//IgcFiles is a slice for storing igc files
var igcFiles []Track

func findIndex(x map[int]string, y int) bool {
	for k, _ := range x {
		if k == y {
			return false
		}
	}
	return true
}

//this function the key of the string if the map contains it, or -1 if the map does not contain the string
func searchMap(x map[int]string, y string) int {

	for k, v := range x {
		if v == y {
			return k
		}
	}
	return -1
}

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
	ID       string    `json:"id"`
	IgcTrack igc.Track `json:"igc_track"`
}

//Attributes : the info about each igc file via id
type Attributes struct {
	HeaderDate string  `json:"h_date"`
	Pilot      string  `json:"pilot"`
	Glider     string  `json:"glider"`
	GliderID   string  `json:"glider_id"`
	Length     float64 `json:"track_length"`
}

//Calculating uptime based on ISO 8601
func timeSince(t time.Time) string {

	Decisecond := 100 * time.Millisecond
	Day := 24 * time.Hour

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
func handler(w http.ResponseWriter, r *http.Request) {

	//Handling for /igcinfo and for /<rubbish>

	http.Error(w, "404 - Page not found!", http.StatusNotFound)

}
func handler1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	//var empty = regexp.MustCompile(``)
	var api = regexp.MustCompile(`api`)

	//Handling for /igcinfo/api
	if len(parts) != 3 || !api.MatchString(parts[2]) {
		http.Error(w, "400 - Bad Request, too many url arguments.", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "{"+"\"uptime\": \""+timeSince(timeStarted)+"\","+"\"info\": \"Service for IGC tracks.\","+"\"version\": \"v1\""+"}")

}

//Handling for /igcinfo/api/igc
func handler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	//Handling GET /igcinfo/api/igc for returning all ids storing in a slice
	case http.MethodGet:
		ids := make([]string, 0)

		for i := range igcFiles {
			ids = append(ids, igcFiles[i].ID)
		}

		json.NewEncoder(w).Encode(ids)

	case http.MethodPost:

		//handling post /igcinfo/api/igc for sending a url and returning an id for that url
		pattern := ".*.igc"

		URL := &_url{}

		var error = json.NewDecoder(r.Body).Decode(URL)
		if error != nil {
			fmt.Fprintln(w, "Error!! ", error)
			return
		}
		rand.Seed(time.Now().UnixNano())

		res, err := regexp.MatchString(pattern, URL.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res {
			mapID = searchMap(urlMap, URL.URL)
			ID := rand.Intn(1000)

			track, _ := igc.ParseLocation(URL.URL)
			if mapID == -1 {
				if findIndex(urlMap, initialID) {
					uniqueId = ID
					urlMap[uniqueId] = URL.URL
					igcFile := Track{}
					igcFile.ID = strconv.Itoa(ID)
					igcFile.IgcTrack = track

					igcFiles = append(igcFiles, igcFile)

					json.NewEncoder(w).Encode(igcFile.ID)
					return
				} else {
					rand.Seed(time.Now().UnixNano())
					uniqueId = rand.Intn(1000)
					urlMap[uniqueId] = URL.URL
					igcFile := Track{}
					igcFile.ID = strconv.Itoa(ID)
					igcFile.IgcTrack = track
					igcFiles = append(igcFiles, igcFile)
					json.NewEncoder(w).Encode(igcFile.ID)
					return

				}
			} else {
				uniqueId = searchMap(urlMap, URL.URL)
				json.NewEncoder(w).Encode(uniqueId)
				return

			}
		}
	default:
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

}
func handler3(w http.ResponseWriter, r *http.Request) {
	//Handling /igcinfo/api/igc/<id>

	w.Header().Set("Content-Type", "application/json")
	idURL := mux.Vars(r)

	rNum, _ := regexp.Compile(`[0-9]+`)
	if !rNum.MatchString(idURL["id"]) {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	attributes := &Attributes{}

	for i := range igcFiles {

		if igcFiles[i].ID == idURL["id"] {
			attributes.HeaderDate = igcFiles[i].IgcTrack.Header.Date.String()
			attributes.Pilot = igcFiles[i].IgcTrack.Pilot
			attributes.Glider = igcFiles[i].IgcTrack.GliderType
			attributes.GliderID = igcFiles[i].IgcTrack.GliderID
			attributes.Length = trackLength(igcFiles[i].IgcTrack)

			json.NewEncoder(w).Encode(attributes)
		}
		//Handling if user type different id from ids stored

	}
	http.Error(w, "404 - The trackInfo with that id doesn't exists in IGC Files", http.StatusNotFound)

}

func handler4(w http.ResponseWriter, r *http.Request) {

	//Handling for GET /api/igc/<id>/<field>
	w.Header().Set("Content-Type", "application/json")

	urlFields := mux.Vars(r)

	var rNum, _ = regexp.Compile(`[a-zA-Z_]+`)

	if !rNum.MatchString(urlFields["field"]) {
		http.Error(w, "400 - Bad Request, wrong parameters", http.StatusBadRequest)
		return
	}

	for i := range igcFiles {

		if igcFiles[i].ID == urlFields["id"] {
			switch {
			case urlFields["field"] == "pilot":
				json.NewEncoder(w).Encode(igcFiles[i].IgcTrack.Pilot)
			case urlFields["field"] == "glider":
				json.NewEncoder(w).Encode(igcFiles[i].IgcTrack.GliderType)

			case urlFields["field"] == "glider_id":
				json.NewEncoder(w).Encode(igcFiles[i].IgcTrack.GliderID)

			case urlFields["field"] == "track_length":
				json.NewEncoder(w).Encode(trackLength(igcFiles[i].IgcTrack))

			case urlFields["field"] == "h_date":
				json.NewEncoder(w).Encode(igcFiles[i].IgcTrack.Header.Date.String())

			default:
				http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
				return
			}

		} else {
			http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
			return
		}

	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/igcinfo", handler)
	r.HandleFunc("/igcinfo/api", handler1)
	r.HandleFunc("/igcinfo/api/igc", handler2)
	r.HandleFunc("/igcinfo/api/igc/{id}", handler3)
	r.HandleFunc("/igcinfo/api/igc/{id}/{field}", handler4)

	//fmt.Println("listening...")

	/*err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}*/
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
