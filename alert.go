package thehive

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"strings"
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
	Date        string     `json:"date,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	Artifacts   []Artifact `json:"artifacts"`
	Raw         []byte     `json:"-"`
}

type AlertResponse struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    int        `json:"severity"`
	Tlp         int        `json:"tlp"`
	Tags        []string   `json:"tags"`
	Type        string     `json:"type"`
	Id          string     `json:"id"`
	Id_         string     `json:"_id"`
	Source      string     `json:"source"`
	SourceRef   string     `json:"sourceRef"`
	Owner       string     `json:"owner"`
	Artifacts   []Artifact `json:"artifacts"`
	Raw         []byte     `json:"-"`
}

// Stores multiple alerts from searches
type HiveAlertMulti struct {
	Raw    []byte          `json:"-"`
	Detail []AlertResponse `json:"-"`
}

// Helperfunction to create an alertartifact based on input
// FIX - does not work for fileupload currently - use struct
func AlertArtifact(dataType string, message string, tlp int, tags []string, ioc bool) Artifact {
	var curartifact Artifact

	// This is weird :)
	curartifact = Artifact{
		DataType: dataType,
		Data:     message,
		Message:  message,
		Tlp:      tlp,
		Tags:     tags,
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
	err = json.Unmarshal(ret.Bytes(), &parsedRet.Detail)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines the creation of an alert
// Takes 10 parameters:
//  1. artifacts []Artifact
//  2. title string
//  3. description string
//  4. tlp int
// 	5. severity int
// 	6. tags []string
//  7. types string
// 	8. source string
// 	9. sourceref string
// 	9. date string
// Returns HiveAlert struct and response error
func (hive *Hivedata) CreateAlert(artifacts []Artifact, title string, description string, tlp int, severity int, tags []string, types string, source string, sourceref string, date string) (*AlertResponse, error) {

	var alert HiveAlert
	var url string

	// Handle files
	newArtifacts := []Artifact{}
	for _, item := range artifacts {
		if item.DataType != "file" {
			newArtifacts = append(newArtifacts, item)
		}

		fd, err := grequests.FileUploadFromDisk(item.Data)
		if err != nil {
			fmt.Println("here?")
			continue
		}

		if fd[0].FileMime == "" {
			fd[0].FileMime = "text/plain"
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(fd[0].FileContents)

		realData := base64.StdEncoding.EncodeToString([]byte(buf.String()))
		if err != nil {
			fmt.Println(err)
			continue
		}

		filenamesplit := strings.Split(fd[0].FileName, "/")
		filename := filenamesplit[(len(filenamesplit))-1]

		item.Data = fmt.Sprintf("%s;%s;%s", filename, fd[0].FileMime, realData)
		newArtifacts = append(newArtifacts, item)
	}

	alert = HiveAlert{
		Title:       title,
		Description: description,
		Tlp:         tlp,
		Artifacts:   newArtifacts,
		Type:        types,
		Tags:        tags,
		SourceRef:   sourceref,
		Source:      source,
		Severity:    severity,
	}

	if date != "" {
		alert.Date = date
	}

	jsondata, err := json.Marshal(alert)

	if err != nil {
		return &AlertResponse{}, err
	}

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url = fmt.Sprintf("%s%s", hive.Url, "/api/alert")
	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(AlertResponse)
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
func (hive *Hivedata) PatchAlertFieldString(alertId string, field string, value string) (*AlertResponse, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)

	data := fmt.Sprintf(`{"%s": "%s"}`, field, value)
	jsondata := []byte(data)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(AlertResponse)
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
func (hive *Hivedata) PatchAlertFieldInt(alertId string, field string, value int) (*AlertResponse, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)

	data := fmt.Sprintf(`{"%s": %d}`, field, value)
	jsondata := []byte(data)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(AlertResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines the modification of artifacts in an alert
// Takes two parameters:
//  1. alertId string
//  2. value []Artifact
// Returns HiveAlert struct and response error
func (hive *Hivedata) AddAlertArtifact(alertId string, artifact Artifact) (*AlertResponse, error) {
	var ret *grequests.Response

	resp, err := hive.GetAlert(alertId)
	if err != nil {
		return nil, err
	}

	if artifact.DataType != "file" {
		resp.Artifacts = append(resp.Artifacts, artifact)

		// Custom solution because files suck. Couldn't get it to work with hive.Ro for some reason
		marshalData, err := json.Marshal(resp.Artifacts)
		if err != nil {
			return nil, err
		}

		//data := map[string]string{"artifacts": fmt.Sprintf("[%s]", string(marshalData))}
		baseData := fmt.Sprintf("%s", string(marshalData))
		data := fmt.Sprintf(`{"artifacts": %s}`, baseData)

		// Set fieldname for every single one

		requestOptions := &grequests.RequestOptions{
			//Files: fd,
			//Data: data,
			RequestBody: bytes.NewReader([]byte(data)),
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
				"Content-Type":  "application/json",
			},
			InsecureSkipVerify: !false,
		}

		url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)
		ret, err = grequests.Patch(url, requestOptions)
	} else {
		fd, err := grequests.FileUploadFromDisk(artifact.Data)
		if err != nil {
			fmt.Println("here?")
			return nil, err
		}

		artifact.Data = ""
		fd[0].FieldName = "attachment"

		resp.Artifacts = []Artifact{}
		resp.Artifacts = append(resp.Artifacts, artifact)

		// Custom solution because files suck. Couldn't get it to work with hive.Ro for some reason
		marshalData, err := json.Marshal(resp.Artifacts)
		if err != nil {
			fmt.Println("HERE)_))")
			return nil, err
		}

		// FIXME - find more files etc
		baseData := fmt.Sprintf("%s", string(marshalData))
		tmpData := fmt.Sprintf(`{"artifacts": %s}`, baseData)
		//data := map[string]string{"artifacts": fmt.Sprintf("%s", tmpData)}
		fmt.Println(tmpData)

		requestOptions := &grequests.RequestOptions{
			Files: fd,
			//Data:  data,
			RequestBody: bytes.NewReader([]byte(tmpData)),
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
				"Content-Type":  "application/json",
			},
			InsecureSkipVerify: !false,
		}

		url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)
		fmt.Println(url)
		ret, err = grequests.Patch(url, requestOptions)
	}

	parsedRet := new(AlertResponse)
	err = json.Unmarshal(ret.Bytes(), parsedRet)
	if err != nil {
		return nil, err
	}

	parsedRet.Raw = ret.Bytes()

	return parsedRet, nil
}

// Removes current tags and adds new ones
// Takes two parameters:
//  1. alertId string
//  2. value []string
// Returns HiveAlert struct and response error
func (hive *Hivedata) PatchAlertTags(alertId string, value []string) (*AlertResponse, error) {
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

	parsedRet := new(AlertResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

func (hive *Hivedata) MarkAlertAsUnread(alertId string) (*AlertResponse, error) {
	url := fmt.Sprintf("%s/api/alert/%s/markAsUnread", hive.Url, alertId)
	ret, err := grequests.Post(url, &hive.Ro)
	if err != nil {
		return nil, err
	}

	parsedRet := new(AlertResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err

}

func (hive *Hivedata) MarkAlertAsRead(alertId string) (*AlertResponse, error) {
	url := fmt.Sprintf("%s/api/alert/%s/markAsRead", hive.Url, alertId)
	ret, err := grequests.Post(url, &hive.Ro)
	if err != nil {
		return nil, err
	}

	parsedRet := new(AlertResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err

}

func (hive *Hivedata) GetAlert(alertId string) (*AlertResponse, error) {
	url := fmt.Sprintf("%s/api/alert/%s", hive.Url, alertId)
	ret, err := grequests.Get(url, &hive.Ro)
	if err != nil {
		return nil, err
	}

	parsedRet := new(AlertResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err

}
