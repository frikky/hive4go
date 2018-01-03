package thehive

/*
	"TheHive" golang api.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

// Stores login data
type Hivedata struct {
	Url      string
	Username string
	Password string
	Ro       grequests.RequestOptions
}

// Stores a hive case
type Hivecase struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Tlp         int        `json:"tlp"`
	Severity    int        `json:"severity"`
	Tags        []string   `json:"tags"`
	Tasks       []CaseTask `json:"tasks"`
	Flag        bool       `json:"flag"`
}

// Stores an artifact
type Artifact struct {
	DataType string   `json:"dataType"`
	Data     string   `json:"data"`
	Tlp      int      `json:"tlp"`
	Tags     []string `json:"tags"`
	Ioc      bool     `json:"ioc"`
}

// Stores alertdata
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

type CaseTask struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Flag        bool   `json:"flag"`
}

// FIX - missing file upload
type CaseTaskLog struct {
	Message string `json:"message"`
}

// Defines API login principles that can be reused in requests
func CreateLogin(inurl string, apikey string) Hivedata {
	formattedApikey := fmt.Sprintf("Bearer %s", apikey)
	return Hivedata{
		Url: inurl,
		Ro: grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": formattedApikey,
			},
			RequestTimeout:     time.Duration(10) * time.Second,
			InsecureSkipVerify: true,
		},
	}
}

// FIX - Missing file upload - See Hive4py api.py and models.py
func (hive *Hivedata) CreateTaskLog(taskId string, taskLog CaseTaskLog) (*grequests.Response, error) {
	var url string
	var err error
	var jsondata []byte

	jsondata, err = json.Marshal(taskLog)

	if err != nil {
		return nil, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s/api/case/task/%s/log", hive.Url, taskId)
	ret, err := grequests.Post(url, &hive.Ro)
	return ret, err
}

func (hive *Hivedata) CreateCaseTask(caseId string, casetask CaseTask) (*grequests.Response, error) {
	var url string
	var err error
	var jsondata []byte

	jsondata, err = json.Marshal(casetask)

	if err != nil {
		return nil, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s/api/case/%s/task", hive.Url, caseId)
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}

// Creates a case and returns based on input data
// Missing date
// FIX - All exits
func (hive *Hivedata) CreateCase(title string, description string, tlp int, severity int, tasks []CaseTask, tags []string, flag bool) (*grequests.Response, error) {
	var curcase Hivecase
	var url string

	if title == "" {
		fmt.Println("Title not set.")
		// WHat do I do here? idk
	}

	if description == "" {
		fmt.Println("Description not set.")
	}

	// Creates case struct for json usage
	curcase = Hivecase{
		Title:       title,
		Description: description,
		Tlp:         tlp,
		Severity:    severity,
		Tags:        tags,
		Tasks:       tasks,
		Flag:        flag,
	}

	// Encodes struct as json
	jsondata, err := json.Marshal(curcase)

	if err != nil {
		return nil, err
	}

	// FIX - might point to same memory, so make a duplicate without editing
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/case")
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}

// Creates an alertartifact based on input
func AlertArtifact(dataType string, message string, tlp int, tags []string, ioc bool) Artifact {
	var curartifact Artifact

	curartifact = Artifact{
		DataType: dataType,
		Data:     message,
		Tlp:      tlp,
		Tags:     tags,
		Ioc:      ioc,
	}

	return curartifact
}

// Gets a single case based on ID
func (hive *Hivedata) GetCase(case_id string) (*grequests.Response, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/case/%s", case_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	resp, err := grequests.Get(url, &hive.Ro)
	return resp, err
}

// Finds all cases based on search parameter
func (hive *Hivedata) FindCases(search []byte) (*grequests.Response, error) {
	var url = fmt.Sprintf("%s%s", hive.Url, "/api/case/_search?sort=%2Btlp&range=all")
	hive.Ro.JSON = search

	resp, err := grequests.Post(url, &hive.Ro)
	return resp, err
}

// Gets an alert based on the alert_id
func (hive *Hivedata) GetAlert(alert_id string) (*grequests.Response, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/alert/%s", alert_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	resp, err := grequests.Get(url, &hive.Ro)
	return resp, err
}

// Gets a field and values in the field
func (hive *Hivedata) FindAlertsQuery(queryfield string, queryvalues []string) (*grequests.Response, error) {
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
func (hive *Hivedata) FindAlertsRaw(search []byte) (*grequests.Response, error) {
	var url string
	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert/_search?range=all")

	hive.Ro.JSON = search

	resp, err := grequests.Post(url, &hive.Ro)
	return resp, err
}

// Creates a case
func (hive *Hivedata) CreateAlert(artifacts []Artifact, title string, description string, tlp int, severity int, tags []string, types string, source string, sourceref string) (*grequests.Response, error) {

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
		return nil, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert")
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}

// Runs analysis on an artifact
// FIX - doesn't work yet as it only contacts Cortex
/*
func (hive *Hivedata) AnalyzeArtifact(analyzerName string) (*grequests.Response, error) {
	cortexUrl := "http://127.0.0.1:9001"
	type AnalyzerAttribute struct {
		Tlp      int    `json:"tlp"`
		DataType string `json:"dataType"`
	}

	type AnalyzerData struct {
		Data       string            `json:"data"`
		Attributes AnalyzerAttribute `json:"attributes"`
	}

	data := AnalyzerData{
		Data: "8.8.8.8",
		Attributes: AnalyzerAttribute{
			Tlp:      2,
			DataType: "ip",
		},
	}

	jsondata, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	//url := fmt.Sprintf("%s/api/analyzer/%s/run", hive.Url, analyzerName)
	url := fmt.Sprintf("%s/api/analyzer/%s/run", cortexUrl, analyzerName)
	ret, err := grequests.Post(url, &hive.Ro)

	return ret, err
}
*/
