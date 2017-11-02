# Hive4go
Statically typed TheHive API for Golang. 

# Missing
Since it's not released yet, just clone to github.com/ in $GOPATH
```Go
import "github.com/frikky/hive4go
```

# Usage
Set login
```Go
login = thehive.CreateLogin("ip", "username", "password")
```

Create Case
```Go
tlp, severity := 3
resp, err := thehive.CreateCase(login, "hive4go title", "hive4go desc", tlp, severity, []string{"task"}, []string{"tags"})
```

# Todo (In order~)
* [FINISHED] Missing all the gets, got all the posts
* [FINISHED] Added most of the get methods
* Create working case POST search (Copy alert)
* Return raw json and not grequests.response (?) (Missing resp.String() formatted to JSON)
* [ALMOST FINISHED] Written tests, needs formatting for publishing 
* Requirements file for running (e.g. grequests)
* Thorougly test get and post methods. Queries don't work properly yet.
* Create an actual readme
* Add to \"go get\" repothingy 
* Implement kwargs somehow (Currently statically typed, keep old stuff too)
