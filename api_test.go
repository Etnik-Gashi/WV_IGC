package WV_IGC

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"server.go"
)


func Test_handler2_notImplemented(t *testing.T) {
	// instantiate mock HTTP server
	// register our handlerStudent <-- actual logic
	ts := httptest.NewServer(http.HandlerFunc(handler2))
	defer ts.Close()

	// create a request to our mock HTTP server
	//    in our case it means to create DELETE request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	// check if the response from the handler is what we expect
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusNotImplemented, resp.StatusCode)
	}
}

func Test_handler2_malformedURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler2))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/igcinfo/api/extra",
		ts.URL + "/igc/",
	}
	for _, tstring := range testCases {
		resp, err := http.Get(tstring)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d", tstring,
				http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}
func Test_handler2_getAllIds_empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler2))
	defer ts.Close()


	resp, err := http.Get(ts.URL + "/igc/")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 0 {
		t.Errorf("Excpected empty array, got %s", a)
	}
}

func Test_handler3_getAllById_59(t *testing.T) {

	testTrack := Track{"59"}

	ts := httptest.NewServer(http.HandlerFunc(HandlerStudent))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/student/")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []Student
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 1 {
		t.Errorf("Excpected array with one element, got %v", a)
	}

	if a[0].Id != testTrack.Id || a[0].igcTrack != testTrack.igcTrack  {
		t.Errorf("Students do not match! Got: %v, Expected: %v\n", a[0], testTrack)
	}
}
