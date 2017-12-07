# Hive4go
Statically typed TheHive API for Golang.   

Based on https://github.com/CERT-BDF/TheHive4py

# Install
```Bash
go get github.com/frikky/hive4go
```

```Go
import "github.com/frikky/hive4go
```

# Usage
Set login, used as first parameter to all functions
```Go
login := thehive.CreateLogin("ip", "apikey")
```

Create case example
```Go
tlp, severity := 3
resp, err := thehive.CreateCase(login, "hive4go title", "hive4go desc", tlp, severity, []string{"task"}, []string{"tags"})
```

# Todo (In order~)
* [FINISHED] Missing all the gets, got all the posts
* [FINISHED] Added most of the get methods
* [FINISHED] Create working case POST search (Copy alert)
* [FINISHED] Thorougly test get and post methods. Queries don't work properly yet.
* [FINISHED] Deprecate BasicAuth (2.12+)
* [ALMOST FINISHED] Written tests, needs formatting for publishing 
* Return raw json and not grequests.response (?) (Missing resp.String() or resp.Bytes())
* Implement kwargs somehow (Currently statically typed, keep old stuff too)
* Requirements file for running (e.g. grequests)
* Create an actual readme
