package thehive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
)

// Stores a hive case
type HiveCase struct {
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Tlp          int                    `json:"tlp"`
	Severity     int                    `json:"severity"`
	Tags         []string               `json:"tags"`
	Tasks        []CaseTask             `json:"tasks"`
	Flag         bool                   `json:"flag"`
	Date         int64                  `json:"date,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Id           string                 `json:"id,omitempty"`
	Owner        string                 `json:"owner,omitempty"`
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
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Tlp              int                    `json:"tlp"`
	Severity         int                    `json:"severity"`
	Date             int64                  `json:"date,omitempty"`
	Tags             []string               `json:"tags"`
	Tasks            []CaseTask             `json:"tasks"`
	Flag             bool                   `json:"flag"`
	Owner            string                 `json:"owner"`
	Status           string                 `json:"status"`
	CreatedAt        int64                  `json:"createdAt"`
	CustomFields     map[string]interface{} `json:"customFields"`
	Id               string                 `json:"id"`
	Summary          string                 `json:"summary"`
	ResolutionStatus string                 `json:"resolutionStatus"`
	ImpactStatus     string                 `json:"impactStatus"`
	Raw              []byte                 `json:"-"`
}

// Stores a response if there are multiple cases
type HiveCaseRespMulti struct {
	Raw    []byte
	Detail []HiveCaseResp
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
func (hive *Hivedata) GetCase(case_id string) (*HiveCaseResp, error) {
	var url, urlpath string

	urlpath = fmt.Sprintf("/api/case/%s", case_id)
	url = fmt.Sprintf("%s%s", hive.Url, urlpath)

	ret, err := grequests.Get(url, &hive.Ro)

	parsedRet := new(HiveCaseResp)
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

func (hive *Hivedata) PatchCaseFieldInt(alertId string, field string, value int64) (*HiveCase, error) {
	url := fmt.Sprintf("%s/api/case/%s", hive.Url, alertId)

	data := fmt.Sprintf(`{"%s": %d}`, field, value)
	fmt.Println(data)
	jsondata := []byte(data)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Patch(url, &hive.Ro)

	parsedRet := new(HiveCase)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}
