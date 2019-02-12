package thehive

import (
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
			RequestTimeout:     time.Duration(30) * time.Second,
			InsecureSkipVerify: !verify,
		},
	}
}
