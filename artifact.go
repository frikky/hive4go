package thehive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
)

// Stores a file attachment response
type FileAttachment struct {
	Name        string   `json:"name"`
	Hashes      []string `json:"hashes"`
	Size        int      `json:"size"`
	ContentType string   `json:"contentType"`
	Id          string   `json:"id"`
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
