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

# Example usage
Set login, used as first parameter to all functions
```Go
login := thehive.CreateLogin("ip", "apikey")
```

Create case example
```Go
resp, err := thehive.CreateCase(
	login,									// Api Login
	"hive4go title", 						// Title
	"hive4go desc", 						// Description
	3, 										// TLP
	3, 										// Severity
	[]thehive.CaseTask{{Title: "task"}}, 	// Tasks
	[]string{"tags"},						// Tags
	true,									// Flag
)
```

# Todos
* [Some finished] Written tests, needs formatting for publishing 
* Return raw json and not grequests.response (?) (Missing resp.String() or resp.Bytes())
* Implement kwargs somehow (Currently statically typed, keep old stuff too)
* Requirements file for running (e.g. grequests)
* Create an actual readme

## Small fixes:)
* Implement files and file add examples
* Implement proxy configuration 
* Implement custom case fields 
* Implement startdate for casetask 
* Implement range and sort 
* Make use of the new search format
