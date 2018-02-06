package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/frikky/hive4go"
	"github.com/satori/go.uuid"
	"os"
)

func main() {
	hive := thehive.CreateLogin("https://192.168.159.137:9443", "lViAsT2BftaFsIN0aBTZm3Ei8EAEp/7P", false)

	// Missing file
	artifacts := []thehive.Artifact{
		thehive.AlertArtifact("ip", "8.8.8.8", 0, []string{}, false),
		thehive.AlertArtifact("domain", "google.com", 0, []string{}, false),
		//thehive.AlertArtifact("file", "pic.png", 0, []string{}, 0)
		//thehive.AlertArtifact("file", "sample.txt", 0, []string{}, 0)
	}

	sourceRef := uuid.NewV4().String()
	fmt.Println("Create Alert")
	fmt.Println("--------------------------")
	alert, err := hive.CreateAlert(
		artifacts,     // Artifacts
		"ASJDLASKLDJ", // Title
		"HELO",        // Description
		1,             // TLP
		1,             // Severity
		[]string{"hive4go", "sample"}, // Tags
		"SIEM",         // Type
		"Carbon black", // Source
		sourceRef,      // SourceRef
	)

	jsonData, err := simplejson.NewJson(alert.Raw)
	ret, _ := jsonData.EncodePretty()
	fmt.Println(string(ret))

	id := jsonData.Get("id").MustString()

	// Get all the details of the created alert
	fmt.Printf("Get created alert %s\n", id)
	fmt.Println("--------------------------")

	response, err := hive.GetAlert(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonData, err = simplejson.NewJson(response.Raw)
	ret, _ = jsonData.EncodePretty()
	fmt.Println(string(ret))
}
