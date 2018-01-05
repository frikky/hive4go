package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/frikky/hive4go"
	"os"
)

func main() {
	hive := thehive.CreateLogin("http://127.0.0.1:9000", "hivekey", false)

	// Example query - Turn into the new search format
	query := `{"query": {"_in": {"_field": "tlp", "_values": [2]}}}`
	queryBytes := []byte(query)
	response, err := hive.FindCases(queryBytes)

	if err != nil || response.StatusCode != 200 {
		fmt.Println(err, response.StatusCode)
		fmt.Println(response.String())
		os.Exit(1)
	}

	jsonData, err := simplejson.NewJson(response.Bytes())
	ret, _ := jsonData.EncodePretty()
	fmt.Println(string(ret))

}
