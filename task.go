package thehive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
)

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
