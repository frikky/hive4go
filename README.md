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

...Alerts etc

# Todo
[FINISHED] Missing all the gets, got all the posts<br>
[FINISHED] Added most of the get methods<br>
[ALMOST FINISHED] Written tests, needs formatting -> publish<br>
Return raw json and not \*grequests.response<br>
Requirements for running (e.g. grequests)<br>
Thorougly test get methods. Queries don't work properly yet.<br>
Create an actual readme<br>
Add to \"go get\" repo<br>
Implement kwargs somehow<br>
