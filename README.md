# Hive4go
Hive4go is a _unofficial_ Golang API client for [TheHive](https://thehive-project.org/).

Based on https://github.com/CERT-BDF/TheHive4py


# Install
```Bash
go get github.com/frikky/hive4go
```

```Go
import "github.com/frikky/hive4go"
```

# Create case example
Set logindata, used for any interactive APIcall 
```Go
verifyCert := false
login := thehive.CreateLogin("ip", "apikey", verifyCert)
```

Create case example
```Go
TLP, Severity := 3
flag := true
resp, err := login.CreateCase(
	"hive4go title", 						
	"hive4go desc", 						
	TLP, 									
	Severity, 								
	[]thehive.CaseTask{{Title: "task"}}, 	
	[]string{"tags"},						
	flag,									
)
```

This will return a case with the following structure. 
```Go
type HiveCase struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Tlp         int        `json:"tlp"`
	Severity    int        `json:"severity"`
	Tags        []string   `json:"tags"`
	Tasks       []CaseTask `json:"tasks"`
	Flag        bool       `json:"flag"`
	Raw         []byte     `json:"-"`
}
```

All return types (alerts, artifacts etc.) follow this type. If you want to handle 
it as raw json, use response.Raw.

# Todos
* [Some finished] Write tests for functions
* Create an actual readme
* Push function documentation
* Implement fileupload
* Implement proxy configuration 
* Implement custom case fields 
* Implement datestamps 
* Make use of the new search format (query.py in hive4py)
* Cleanup
