package main

import (
	"fmt"
	"github.com/frikky/hive4go"
)

func main() {
	hive := thehive.CreateLogin("http://127.0.0.1:9000", "apikey", false)

	ret, err := hive.GetTaskLogs("Case-task ID")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ret)
}
