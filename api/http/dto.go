package http

import (
	"encoding/json"
	"fraisedb/base"
	"io"
	"net/http"
)

type nodeCommand struct {
	Addr     string `json:"addr"`
	TcpPort  int    `json:"tcpPort"`
	HttpPort int    `json:"httpPort"`
}

type kvCommand struct {
	SaveType int    `json:"type"`
	Value    string `json:"value"`
	Incr     int64  `json:"incr"`
	TTL      int64  `json:"ttl"`
}

func reply(w http.ResponseWriter, result any, err error) {
	var res = make(map[string]any, 2)
	if err != nil {
		res["error"] = err.Error()
		base.LogHandler.Println(base.LogErrorTag, err)
	}
	if result != nil {
		res["result"] = result
	}
	marshal, _ := json.Marshal(res)
	_, err = w.Write(marshal)
	if err != nil {
		base.LogHandler.Println(base.LogErrorTag, err)
	}
}

func readKVCommand(r *http.Request) (kvCommand, error) {
	command := kvCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return command, err
	}
	err = json.Unmarshal(body, &command)
	return command, err
}

func readNodeCommand(r *http.Request) (nodeCommand, error) {
	command := nodeCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return command, err
	}
	err = json.Unmarshal(body, &command)
	return command, err
}
