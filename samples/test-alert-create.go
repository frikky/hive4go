package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/frikky/hive4go"
	"github.com/satori/go.uuid"
	"os"
)

func main() {
	hive := thehive.CreateLogin("http://127.0.0.1:9000", "apikey")

	// Missing file
	artifacts := []thehive.Artifact{
		hive.AlertArtifact("ip", "8.8.8.8", 0, []string{}, false),
		hive.AlertArtifact("domain", "google.com", 0, []string{}, false),
		//thehive.AlertArtifact("file", "pic.png", 0, []string{}, 0)
		//thehive.AlertArtifact("file", "sample.txt", 0, []string{}, 0)
	}

	sourceRef := uuid.NewV4().String()
	fmt.Println("Create Alert")
	fmt.Println("--------------------------")
	alert, err := hive.CreateAlert(
		artifacts,   // Artifacts
		"New Alert", // Title
		"N/A",       // Description
		3,           // TLP
		1,           // Severity
		[]string{"hive4go", "sample"}, // Tags
		"external",                    // Type
		"instance1",                   // Source
		sourceRef,                     // SourceRef
	)

	if err != nil || alert.StatusCode != 201 {
		fmt.Println(err, alert.StatusCode)
		os.Exit(1)
	}

	jsonData, err := simplejson.NewJson(alert.Bytes())
	ret, _ := jsonData.EncodePretty()
	fmt.Println(string(ret))

	id := jsonData.Get("id").MustString()

	// Get all the details of the created alert
	fmt.Printf("Get created alert %s\n", id)
	fmt.Println("--------------------------")

	response, err := hive.GetAlert(id)
	fmt.Println(response.StatusCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonData, err = simplejson.NewJson(response.Bytes())
	ret, _ = jsonData.EncodePretty()
	fmt.Println(string(ret))
}
