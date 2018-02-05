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
	Url    string
	Apikey string
	Ro     grequests.RequestOptions
}

// Stores a hive case
type HiveCase struct {
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Tlp          int                    `json:"tlp"`
	Severity     int                    `json:"severity"`
	Tags         []string               `json:"tags"`
	Tasks        []CaseTask             `json:"tasks"`
	Flag         bool                   `json:"flag"`
	CustomFields map[string]interface{} `json:"customFields"`
	Raw          []byte                 `json:"-"`
}

// Stores multiple hive cases from searches
type HiveCaseMulti struct {
	Raw    []byte
	Detail []HiveCase
}

// Stores the response of a case from thehive
type HiveCaseResp struct {
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Tlp          int                    `json:"tlp"`
	Severity     int                    `json:"severity"`
	Tags         []string               `json:"tags"`
	Tasks        []CaseTask             `json:"tasks"`
	Flag         bool                   `json:"flag"`
	CreatedAt    int64                  `json:"createdAt"`
	CustomFields map[string]interface{} `json:"customFields"`
	Id           string                 `json:"id"`
	Raw          []byte                 `json:"-"`
}

// Stores a response if there are multiple cases
type HiveCaseRespMulti struct {
	Raw    []byte
	Detail []HiveCaseResp
}

// Stores an artifact
type Artifact struct {
	DataType string                 `json:"dataType"`
	Message  string                 `json:"message"`
	Tlp      int                    `json:"tlp"`
	Tags     []string               `json:"tags"`
	Ioc      bool                   `json:"ioc"`
	Files    []grequests.FileUpload `json:"-"`
	Data     string                 `json:"-"`
	Raw      []byte                 `json:"-"`
}

// Stores multiple artifacts from searches
type HiveArtifactMulti struct {
	Raw    []byte
	Detail []Artifact
}

// Missing Reports
type ArtifactResponse struct {
	DataType   string         `json:"dataType"`
	CreatedBy  string         `json:"createdBy"`
	Sighted    bool           `json:"sighted"`
	Tlp        int            `json:"tlp"`
	_Id        string         `json:"_id"`
	Tags       []string       `json:"tags"`
	Message    string         `json:"message"`
	Ioc        bool           `json:"ioc"`
	Status     string         `json:"status"`
	Attachment FileAttachment `json:"attachment"`
	Id         string         `json:"id"`
	Type       string         `json:"_type"`
	Raw        []byte         `json:"-"`
}

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

// Stores a hive task
type CaseTask struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Flag        bool   `json:"flag"`
	Raw         []byte `json:"-"`
}

// Stores multiple tasks from searches
type CaseTaskMulti struct {
	Raw    []byte
	Detail []CaseTask
}

// Stores task responses
type CaseTaskResponse struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Flag        bool   `json:"flag"`
	CreatedBy   string `json:"createdBy"`
	Order       int    `json:"order"`
	Id          string `json:"id"`
	Type        string `json:"_type"`
	Raw         []byte `json:"-"`
}

// Stores multiple task responses  from searches
type CaseTaskRespMulti struct {
	Raw    []byte
	Detail []CaseTaskResponse
}

// Should maybe have title etc as well.
type CaseTaskLog struct {
	Message string                 `json:"message"`
	Files   []grequests.FileUpload `json:"-"`
	Raw     []byte                 `json:"-"`
}

// Stores multiple casetasklogs
type CaseTaskLogMulti struct {
	Raw    []byte
	Detail []CaseTaskLog
}

// Stores a file attachment response
type FileAttachment struct {
	Name        string   `json:"name"`
	Hashes      []string `json:"hashes"`
	Size        int      `json:"size"`
	ContentType string   `json:"contentType"`
	Id          string   `json:"id"`
}

// Stores case tasklog responses
type CaseTaskLogResponse struct {
	Message    string         `json:"message"`
	Title      string         `json:"title"`
	CreatedBy  string         `json:"createdBy"`
	Order      string         `json:"order"`
	Owner      string         `json:"owner"`
	Flag       bool           `json:"flag"`
	Status     string         `json:"status"`
	Id         string         `json:"id"`
	Type       string         `json:"_type"`
	Attachment FileAttachment `json:"attachment"`
	Raw        []byte         `json:"-"`
}

// Stores multiple tasklog responses
type CaseTaskLogRespMulti struct {
	Raw    []byte
	Detail []CaseTaskLogResponse
}

// Defines API login principles that can be reused in requests
// Takes three parameters:
//  1. URL string
//  2. API key
//  3. Verify boolean that should be true in order to verify the servers certificate
// Returns Hivedata struct
func CreateLogin(inurl string, apikey string, verify bool) Hivedata {
	formattedApikey := fmt.Sprintf("Bearer %s", apikey)
	return Hivedata{
		Url:    inurl,
		Apikey: apikey,
		Ro: grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": formattedApikey,
			},
			RequestTimeout:     time.Duration(10) * time.Second,
			InsecureSkipVerify: !verify,
		},
	}
}

// Defines case task log creation
// Takes two parameters:
//  1. taskId string
//  2. taskLog CaseTaskLog
// Returns CaseTaskLogresponse struct and response error
func (hive *Hivedata) CreateTaskLog(taskId string, taskLog CaseTaskLog) (*CaseTaskLogResponse, error) {
	var url string
	var err error
	var jsondata []byte

	jsondata, err = json.Marshal(taskLog)
	if err != nil {
		return nil, err
	}

	// If files exist
	if taskLog.Files != nil {
		hive.Ro.Files = taskLog.Files
		hive.Ro.Data = map[string]string{
			"_json": string(jsondata),
		}
		hive.Ro.Headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
		}
	} else {
		hive.Ro.RequestBody = bytes.NewReader(jsondata)
	}

	url = fmt.Sprintf("%s/api/case/task/%s/log", hive.Url, taskId)
	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(CaseTaskLogResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines case task log creation
// Takes three parameters:
//  1. caseId string
//  2. name string
//  3. data string
// Returns CaseTaskLogresponse struct and response error
// FIX - only supports string currently
func (hive *Hivedata) AddCustomFieldData(caseId string, name string, data string) (*HiveCase, error) {
	jsonQuery := fmt.Sprintf(`{"customFields.%s": {"string": "%s"}}`, name, data)
	jsondata := []byte(jsonQuery)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url := fmt.Sprintf("%s/api/case/%s", hive.Url, caseId)
	//resp, err := grequests.Post(url, &hive.Ro)
	resp, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(HiveCase)
	_ = json.Unmarshal(resp.Bytes(), parsedRet)
	parsedRet.Raw = resp.Bytes()

	return parsedRet, err
}

// Defines creation of a case task within a case
// Takes two parameters:
//  1. caseId string
//  2. casetask CaseTask
// Returns CaseTask struct and response error
func (hive *Hivedata) CreateCaseTask(caseId string, casetask CaseTask) (*CaseTaskResponse, error) {
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

	parsedRet := new(CaseTaskResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines creation of a case
// Takes two parameters:
//  1. title string
//  2. description string
//  3. tlp int
// 	4. severity int
// 	5. tasks []CaseTask
// 	6. tags []string
// 	7. flag bool
// Returns HiveCase struct and response error
func (hive *Hivedata) CreateCase(title string, description string, tlp int, severity int, tasks []CaseTask, tags []string, flag bool) (*HiveCase, error) {
	var curcase HiveCase
	var url string

	if title == "" {
		fmt.Println("Missing title in API call. Set title.")
		title = "Missing title in API call"
	}

	if description == "" {
		fmt.Println("Description not set.")
		description = ""
	}

	// Creates case struct for json usage
	curcase = HiveCase{
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

	parsedRet := new(HiveCase)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Creates an alertartifact based on input
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

// Defines how to get  a taskwithin a case
// Takes one parameters:
//  1. taskId string
// Returns CaseTaskResponse struct and response error
func (hive *Hivedata) GetTask(taskId string) (*CaseTaskResponse, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/case/task/%s/log", taskId)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Get(url, &hive.Ro)

	parsedRet := new(CaseTaskResponse)
	_ = json.Unmarshal(ret.Bytes(), &parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Defines creation of a case task within a case
// Takes two parameters:
//  1. title string
//  2. description string
//  3. tlp int
// 	4. severity int
// 	5. tasks []CaseTask
// 	6. tags []string
// 	7. flag bool
// Returns HiveCase struct and response error
func (hive *Hivedata) GetCase(case_id string) (*HiveCase, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/case/%s", case_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Get(url, &hive.Ro)

	parsedRet := new(HiveCase)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Finds all cases based on search parameter
// Takes one argument
//  1. search []byte defined as a marshalled json string
// FIX - Missing sort and range, can use struct
func (hive *Hivedata) FindCases(search []byte) (*HiveCaseRespMulti, error) {
	var url = fmt.Sprintf("%s%s", hive.Url, "/api/case/_search?sort=%2Btlp&range=all")
	hive.Ro.JSON = search

	ret, err := grequests.Post(url, &hive.Ro)

	// Not yet fixed
	parsedRet := new(HiveCaseRespMulti)
	_ = json.Unmarshal(ret.Bytes(), &parsedRet.Detail)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Gets an alert based on the alertId
// Takes one argument
// 1. alertId string
// Returns the alert and the response error
func (hive *Hivedata) GetAlert(alertId string) (*HiveAlert, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/alert/%s", alertId)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Get(url, &hive.Ro)

	parsedRet := new(HiveAlert)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Gets all tasks for a specific case
// Takes one argument
// 1. caseId string
// Returns the casetasks and the response error
func (hive *Hivedata) GetCaseTasks(caseId string) (*CaseTaskRespMulti, error) {
	urlpath := fmt.Sprintf("/api/case/%s/task/_search?range=all", caseId)
	jsonQuery := fmt.Sprintf(`{"_parent": {"_type": "case", "_query": {"id": "%s"}}}`, caseId)
	jsondata := []byte(jsonQuery)

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url := fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(CaseTaskRespMulti)
	json.Unmarshal(ret.Bytes(), &parsedRet.Detail)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}

// Gets all tasks for a specific case
// Takes one argument
// 1. taskId string
// Returns the casetasks and the response error
func (hive *Hivedata) GetTaskLogs(taskId string) (*CaseTaskLogRespMulti, error) {
	// Remove the header as the endpoint doesn't accept application/json..
	hive.Ro.Headers = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
	}

	urlpath := fmt.Sprintf("/api/case/task/%s/log?range=all", taskId)
	url := fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Get(url, &hive.Ro)

	parsedRet := new(CaseTaskLogRespMulti)
	parsedRet.Raw = ret.Bytes()
	json.Unmarshal(parsedRet.Raw, &parsedRet.Detail)

	// Rebuild it
	hive.Ro.Headers = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
	}

	return parsedRet, err
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

// FIX - Doesn't map back to a struct yet
func (hive *Hivedata) AnalyzeArtifact(cortexId string, artifactId string, analyzerId string) (*grequests.Response, error) {
	rawJson := fmt.Sprintf(`{"cortexId":"%s","artifactId":"%s","analyzerId":"%s"}`, cortexId, artifactId, analyzerId)
	jsondata := []byte(rawJson)

	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	url := fmt.Sprintf("%s/api/connector/cortex/job", hive.Url)
	ret, err := grequests.Post(url, &hive.Ro)

	if err != nil {
		return ret, err
	}

	return ret, nil
}

// Defines the creation of an artifact
// Takes two parameters:
//  1. caseId string
//  2. caseArtifact Artifact struct
// Returns ArtifactResponse struct and response error
func (hive *Hivedata) AddCaseArtifact(caseId string, caseArtifact Artifact) (*ArtifactResponse, error) {
	url := fmt.Sprintf("%s/api/case/%s/artifact", hive.Url, caseId)
	jsondata, err := json.Marshal(caseArtifact)
	fmt.Println(string(jsondata))

	if err != nil {
		return new(ArtifactResponse), err
	}

	if caseArtifact.DataType == "file" {
		fileToUpload, err := grequests.FileUploadFromDisk(caseArtifact.Data)
		fileToUpload[0].FieldName = "attachment"

		if err != nil {
			return new(ArtifactResponse), err
		}

		hive.Ro.Files = fileToUpload
		hive.Ro.Data = map[string]string{
			"_json": string(jsondata),
		}
		hive.Ro.Headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
		}
	} else {
		hive.Ro.RequestBody = bytes.NewReader(jsondata)
	}

	ret, err := grequests.Post(url, &hive.Ro)

	parsedRet := new(ArtifactResponse)
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
