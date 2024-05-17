package http

import (
	"fmt"
	"fraisedb/base"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// InitRouter 初始化HTTP路由
func InitRouter() error {

	port := fmt.Sprintf(":%v", base.Config().Node.HttpPort)

	r := httprouter.New()

	r.GET("/health", getHealth)
	r.GET("/config", getConfig)

	r.POST("/node", addNode)
	r.DELETE("/node/:endpoint", removeNode)
	r.GET("/nodes", listNode)
	r.GET("/leader", getLeader)

	r.POST("/namespace/:namespace", createNamespace)
	r.GET("/namespaces", listNamespace)
	r.DELETE("/namespace/:namespace", deleteNamespace)

	r.PUT("/kv/:namespace/:key", putKV)
	r.GET("/kv/:namespace/:key", getKV)
	r.DELETE("/kv/:namespace/:key", deleteKV)
	r.GET("/kvs/:namespace/:keyPrefix", listKV)

	r.GET("/subscribe/:namespace/:keyPrefix/:clientId", subscribe)

	initConsumer()
	err := http.ListenAndServe(port, r)
	if err != nil {
		return err
	}
	return nil
}
