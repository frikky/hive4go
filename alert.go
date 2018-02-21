package thehive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
)

// Stores a hive alert
type HiveAlert struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    int        `json:"severity"`
	Tlp         int        `json:"tlp"`
	Tags        []string   `json:"tags"`
	Type        string     `json:"type"`
	Source      string     `json:"source"`
	SourceRef   string     `json:"sourceRef"`
	Artifacts   []Artifact `json:"artifacts"`
	Raw         []byte     `json:"-"`
}

// Stores multiple alerts from searches
type HiveAlertMulti struct {
	Raw    []byte
	Detail []HiveAlert
}

// Helperfunction to create an alertartifact based on input
// FIX - does not work for fileupload currently - use struct
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

// Creates a search based on field and values - FIX, might be deprecated
// Takes two arguments
// 1. queryfield string
// 2. queryvalues []string
// Returns multiple Alerts and the response error
func (hive *Hivedata) FindAlertsQuery(queryfield string, queryvalues []string) (*HiveAlertMulti, error) {
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

	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(HiveAlertMulti)
	_ = json.Unmarshal(ret.Bytes(), parsedRet.Detail)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Gets a raw json query and returns all data
// Takes one parameter:
//  1. search []bytes - Raw marshalled JSON string
// Returns multiple alerts and the request response
func (hive *Hivedata) FindAlertsRaw(search []byte) (*HiveAlertMulti, error) {
	var url string
	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert/_search?range=all")

	hive.Ro.JSON = search

	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(HiveAlertMulti)
	_ = json.Unmarshal(ret.Bytes(), parsedRet.Detail)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines the creation of an alert
// Takes two parameters:
//  1. artifacts []Artifact
//  2. title string
//  3. description string
//  4. tlp int
// 	5. severity int
// 	6. tags []string
//  7. types string
// 	8. source string
// 	9. sourceref string
// Returns HiveAlert struct and response error
func (hive *Hivedata) CreateAlert(artifacts []Artifact, title string, description string, tlp int, severity int, tags []string, types string, source string, sourceref string) (*HiveAlert, error) {

	var alert HiveAlert
	var url string

	// FUCK THIS :)
	fmt.Println(artifacts)
	alert = HiveAlert{
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
		return &HiveAlert{}, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert")
	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(HiveAlert)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines the modification of an alert
// Takes three parameters:
//  1. alertId string
//  2. field struct
//  3. value struct
// Returns HiveAlert struct and response error
func (hive *Hivedata) PatchAlertFieldString(alertId string, field string, value string) (*HiveAlert, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)

	data := fmt.Sprintf(`{"%s": "%s"}`, field, value)
	jsondata := []byte(data)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(HiveAlert)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines the modification of artifacts in an alert
// Takes two parameters:
//  1. alertId string
//  2. value []Artifact
// Returns HiveAlert struct and response error
func (hive *Hivedata) PatchAlertArtifact(alertId string, value []Artifact) (*HiveAlert, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)

	jsonRet, _ := json.Marshal(value)
	jsondata := []byte(fmt.Sprintf(`{"artifacts": %s}`, string(jsonRet)))

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(HiveAlert)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Removes current tags and adds new ones
// Takes two parameters:
//  1. alertId string
//  2. value []string
// Returns HiveAlert struct and response error
func (hive *Hivedata) PatchAlertTags(alertId string, value []string) (*HiveAlert, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)

	// Better than looping and adding to a string
	type tmpjson struct {
		Tags []string `json:"tags"`
	}

	tmpstruct := tmpjson{}
	tmpstruct.Tags = value

	jsondata, _ := json.Marshal(tmpstruct)

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(HiveAlert)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}
