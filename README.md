# Hive4go
Statically typed TheHive API which so far can create cases and alerts. 

# Missing
Since it's not released yet, just add hive4go to ~/go/src/github.com/frikky/hive4go
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
Add to go get repo<br>
Missing all the gets, got all the posts<br>
Implement kwargs somehow<br>
Create an actual readme<br>
