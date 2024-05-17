package main

import (
	"fmt"
	"fraisedb/api/http"
	"fraisedb/base"
	"fraisedb/core"
	"time"
)

func main() {
	err := _base()
	if err != nil {
		fmt.Println("base error:", err)
		time.Sleep(5 * time.Second)
		panic(err)
	}
	err = _app()
	if err != nil {
		base.LogHandler.Println(base.LogErrorTag, err)
	}
}

func _base() error {
	err := base.InitConfig()
	if err != nil {
		return err
	}
	return base.InitLog()
}

func _app() error {
	err := core.InitNode()
	if err != nil {
		return err
	}
	err = core.JoinCluster()
	if err != nil {
		return err
	}
	return http.InitRouter()
}
