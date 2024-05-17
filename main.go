package main

import (
	"fmt"
	"fraisedb/api/http"
	"fraisedb/base"
	"fraisedb/core"
	"time"
)

func main() {
	err := base.InitConfig()
	if err != nil {
		fmt.Println("init log error:", err)
		time.Sleep(5 * time.Second)
		panic(err)
	}
	err = base.InitLog()
	if err != nil {
		fmt.Println("init log error:", err)
		time.Sleep(5 * time.Second)
		panic(err)
	}
	err = core.InitNode()
	if err != nil {
		base.LogHandler.Println(base.LogErrorTag, err)
		return
	}
	err = http.InitRouter()
	if err != nil {
		base.LogHandler.Println(base.LogErrorTag, err)
		return
	}
}
