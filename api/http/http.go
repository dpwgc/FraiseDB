package http

import (
	"fmt"
	"fraisedb/base"
	"fraisedb/core"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func getHealth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, nil, nil)
}

func getConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, base.Config(), nil)
}

func addNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := readNodeCommand(r)
	if err != nil {
		reply(w, nil, err)
		return
	}
	reply(w, nil, core.AddNode(command.Addr, command.TcpPort, command.HttpPort))
}

func removeNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, nil, core.RemoveNode(p.ByName("endpoint")))
}

func getLeader(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, core.GetLeader(), nil)
}

func listNode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, core.ListNode(), nil)
}

func listNamespace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reply(w, core.ListNamespace(), nil)
}

func createNamespace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if isLeader() {
		forwardToLeader(w, r)
		return
	}
	reply(w, nil, core.CreateNamespace(p.ByName("namespace")))
}

func deleteNamespace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if isLeader() {
		forwardToLeader(w, r)
		return
	}
	reply(w, nil, core.DeleteNamespace(p.ByName("namespace")))
}

func putKV(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if isLeader() {
		forwardToLeader(w, r)
		return
	}
	command, err := readKVCommand(r)
	if err != nil {
		reply(w, nil, err)
		return
	}
	var ddl int64 = 0
	if command.TTL > 0 {
		ddl = time.Now().Unix() + command.TTL
	}
	reply(w, nil, core.PutKV(p.ByName("namespace"), p.ByName("key"), command.Overwrite, command.Value, command.Incr, ddl))
}

func getKV(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	value, err := core.GetKV(p.ByName("namespace"), p.ByName("key"))
	if err != nil {
		reply(w, nil, err)
		return
	}
	reply(w, *value, nil)
}

func getNestedField(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query, err := readNestedQuery(r)
	if err != nil {
		reply(w, nil, err)
		return
	}
	value, err := core.GetNestedField(p.ByName("namespace"), p.ByName("key"), query.Path)
	if err != nil {
		reply(w, nil, err)
		return
	}
	reply(w, value, nil)
}

func deleteKV(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if isLeader() {
		forwardToLeader(w, r)
		return
	}
	reply(w, nil, core.DeleteKV(p.ByName("namespace"), p.ByName("key")))
}

func listKV(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		reply(w, nil, err)
		return
	}
	count, err := strconv.ParseInt(r.URL.Query().Get("count"), 10, 64)
	if err != nil {
		reply(w, nil, err)
		return
	}
	kvs, err := core.ListKV(p.ByName("namespace"), p.ByName("keyPrefix"), offset, count)
	if err != nil {
		reply(w, nil, err)
		return
	}
	reply(w, *kvs, nil)
}

func isLeader() bool {
	// 如果该节点是leader，则无需转发请求
	if core.GetLeader() == fmt.Sprintf("%s:%v", base.Config().Server.Addr, base.Config().Server.HttpPort) {
		return true
	} else {
		return false
	}
}

func forwardToLeader(w http.ResponseWriter, r *http.Request) {
	err := base.HttpForward(w, r, fmt.Sprintf("http://%s%s", core.GetLeader(), r.URL.RequestURI()))
	if err != nil {
		base.LogHandler.Println(base.LogErrorTag, err)
	}
}
