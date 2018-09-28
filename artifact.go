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
	Files    []grequests.FileUpload `json:"-"`
	Data     string                 `json:"data,omitempty"`
	Ioc      bool                   `json:"ioc"`
}

//Raw      []byte                 `json:"-"`

// Stores multiple artifacts from searches
type HiveArtifactMulti struct {
	Raw    []byte
	Detail []Artifact
}

type HiveArtifactMultiResponse struct {
	Raw    []byte
	Detail []ArtifactResponse
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
	Data       string         `json:"data"`
	Ioc        bool           `json:"ioc"`
	Status     string         `json:"status"`
	Attachment FileAttachment `json:"attachment"`
	Id         string         `json:"id"`
	Type       string         `json:"_type"`
	Raw        []byte         `json:"-"`
}

func (hive *Hivedata) GetCaseArtifacts(caseId string) (*HiveArtifactMultiResponse, error) {
	url := fmt.Sprintf("%s/api/case/artifact/_search?range=all", hive.Url)

	rawJson := fmt.Sprintf(
		`{"query": {
		"_and": [{
			"_parent": {
				"_type": "case", 
				"_query": {
					"_id": "AWHCBwmJGM-hjLdzubbW"
				}
			}
		}, {
			"status": "Ok"
		}]
	}}`,
	)

	jsondata := []byte(rawJson)
	hive.Ro.RequestBody = bytes.NewReader(jsondata)

	ret, err := grequests.Post(url, &hive.Ro)

	//type HiveArtifactMulti struct {
	parsedRet := new(HiveArtifactMultiResponse)
	err = json.Unmarshal(ret.Bytes(), &parsedRet.Detail)
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
	var err error
	var ret *grequests.Response

	url := fmt.Sprintf("%s/api/case/%s/artifact", hive.Url, caseId)
	jsondata, err := json.Marshal(caseArtifact)

	if err != nil {
		return new(ArtifactResponse), err
	}

	if caseArtifact.DataType == "file" {
		fd, err := grequests.FileUploadFromDisk(caseArtifact.Data)
		if err != nil {
			return nil, err
		}

		caseArtifact.Data = ""

		// Custom solution because files suck. Couldn't get it to work with hive.Ro for some reason
		marshalData, err := json.Marshal(caseArtifact)
		if err != nil {
			return nil, err
		}

		data := map[string]string{"_json": fmt.Sprintf("%s", string(marshalData))}

		fd[0].FieldName = "attachment"
		requestOptions := &grequests.RequestOptions{
			Files: fd,
			Data:  data,
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", hive.Apikey),
			},
			InsecureSkipVerify: !false,
		}

		ret, err = grequests.Post(url, requestOptions)
	} else {
		hive.Ro.RequestBody = bytes.NewReader(jsondata)
		ret, err = grequests.Post(url, &hive.Ro)
	}

	parsedRet := new(ArtifactResponse)
	_ = json.Unmarshal(ret.Bytes(), parsedRet)
	parsedRet.Raw = ret.Bytes()

	return parsedRet, err
}
