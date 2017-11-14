package thehive

/*
	Attempt at making a "TheHive" api for golang.
	Hive4go?
	Static stuff sucks tho
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	//"io/ioutil"
	//"net/http"
	"os"
)

type Hivedata struct {
	Url      string
	Username string
	Password string
	Ro       grequests.RequestOptions
}

type Hivecase struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tlp         int      `json:"tlp"`
	Severity    int      `json:"severity"`
	Tags        []string `json:"tags"`
	Tasks       []string `json:"tasks"`
}

type Artifact struct {
	DataType string   `json:"dataType"`
	Data     string   `json:"data"`
	Tlp      int      `json:"tlp"`
	Tags     []string `json:"tags"`
}

type AlertData struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    int        `json:"severity"`
	Tlp         int        `json:"tlp"`
	Tags        []string   `json:"tags"`
	Type        string     `json:"type"`
	Source      string     `json:"source"`
	SourceRef   string     `json:"sourceRef"`
	Artifacts   []Artifact `json:"artifacts"`
}

// Defines basic login principles that can be reused in requests
func CreateLogin(inurl string, inusername string, inpassword string) Hivedata {
	logindata := Hivedata{
		Url:      inurl,
		Username: inusername,
		Password: inpassword,
		Ro: grequests.RequestOptions{
			Auth: []string{inusername, inpassword},
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		RequestTimeout: time.duration(10) * time.Second,
	}

	return logindata
}

// Creates a case and returns based on input data
//func createCase(hive Hivedata, args ...args.V) {
// Missing date
func CreateCase(hive Hivedata, title string, description string, tlp int, severity int, tasks []string, tags []string) (*grequests.Response, error) {
	var curcase Hivecase
	var url string

	if title == "" {
		fmt.Println("Title not set.")
		os.Exit(3)
		// WHat do I do here? idk
	}

	if description == "" {
		fmt.Println("Description not set.")
		os.Exit(3)
	}

	// Creates case struct for json usage
	curcase = Hivecase{
		Title:       title,
		Description: description,
		Tlp:         tlp,
		Severity:    severity,
		Tags:        tags,
		Tasks:       tasks,
	}

	// Encodes struct as json
	jsondata, err := json.Marshal(curcase)

	if err != nil {
		fmt.Println("Error in json")
		os.Exit(1)
	}

	// FIX - might point to same memory, so make a duplicate without editing
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/case")
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}

func AlertArtifact(dataType string, message string, tlp int, tags []string) Artifact {
	var curartifact Artifact

	curartifact = Artifact{
		DataType: dataType,
		Data:     message,
		Tlp:      tlp,
		Tags:     tags,
	}

	return curartifact
}

func GetCase(hive Hivedata, case_id string) (*grequests.Response, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("api/case/%s", case_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	resp, err := grequests.Get(url, &hive.Ro)
	return resp, err
}

func FindCases(hive Hivedata, search []byte) (*grequests.Response, error) {
	var url = fmt.Sprintf("%s%s", hive.Url, "/api/case/_search?sort=%2Btlp&range=all")
	hive.Ro.JSON = search

	resp, err := grequests.Post(url, &hive.Ro)
	return resp, err
}

func GetAlert(hive Hivedata, alert_id string) (*grequests.Response, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("api/alert/%s", alert_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	resp, err := grequests.Get(url, &hive.Ro)
	return resp, err
}

// Gets a field and values in the field
func FindAlertsQuery(hive Hivedata, queryfield string, queryvalues []string) (*grequests.Response, error) {
	// Sorts by tlp by default
	var url string

	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert/_search?range=all")

	type Search struct {
		Field  string   `json:"_field"`
		Values []string `json:"_values"`
	}

	type In struct {
		Search `json:"_in"`
	}

	// This one isn't documented, but necessary to make the search work.
	type Query struct {
		In `json:"query"`
	}

	// Creates the json struct object
	searchquery := Query{
		In{
			Search{
				Field:  queryfield,
				Values: queryvalues,
			},
		},
	}

	jsonsearch, err := json.Marshal(searchquery)
	if err != nil {
		return nil, err
	}

	hive.Ro.JSON = jsonsearch

	resp, err := grequests.Post(url, &hive.Ro)
	return resp, err
}

// Gets a raw json query and returns all data
func FindAlertsRaw(hive Hivedata, search []byte) (*grequests.Response, error) {
	var url string
	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert/_search?range=all")

	hive.Ro.JSON = search

	resp, err := grequests.Post(url, &hive.Ro)
	return resp, err
}

// Attempts to create an alert

// Attempts to create an alert
func CreateAlert(hive Hivedata, artifacts []Artifact, title string, description string, tlp int, severity int, tags []string, types string, source string, sourceref string) (*grequests.Response, error) {

	var alert AlertData
	var url string

	alert = AlertData{
		Title:       title,
		Description: description,
		Tlp:         tlp,
		Artifacts:   artifacts,
		Type:        types,
		Tags:        tags,
		SourceRef:   sourceref,
		Source:      source,
		Severity:    severity,
	}

	jsondata, err := json.Marshal(alert)

	if err != nil {
		fmt.Println("Error in json")
		os.Exit(1)
	}

	fmt.Println(string(jsondata))
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert")
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}
