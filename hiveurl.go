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

	url = fmt.Sprintf("http://%s%s", hive.Url, "/api/case")
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

/*
// Notes from python to find cases
func getCase(hive Hivedata) {
	//hive.url = url
	//req = self.url + "/api/case/{}".format(case_id)
	url = fmt.Sprintf("http://%s%s", hive.url, "/api/case")
	return grequests.get(req, hive.ro)
}
req = self.url + "/api/case/{}".format(case_id)

        try:
					        except requests.exceptions.RequestException as e:
							            sys.exit("Error: {}".format(e))
*/

// Attempts to create an alert
func CreateAlert(hive Hivedata, artifacts []Artifact, title string, description string, tlp int, tags []string, types string) (*grequests.Response, error) {

	var alert AlertData
	var url string

	alert = AlertData{
		Title:       title,
		Description: description,
		Tlp:         tlp,
		Artifacts:   artifacts,
		Type:        types,
		Tags:        tags,
		SourceRef:   "#ASD1024",
	}

	jsondata, err := json.Marshal(alert)

	if err != nil {
		fmt.Println("Error in json")
		os.Exit(1)
	}

	fmt.Println(string(jsondata))
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("http://%s%s", hive.Url, "/api/alert")
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}
