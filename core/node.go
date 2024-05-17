package core

import (
	"encoding/json"
	"fmt"
	"fraisedb/base"
	"fraisedb/cluster"
	"fraisedb/store"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// InitNode 启动节点
func InitNode() error {
	base.Channel = make(chan []byte, 1000)
	err := base.CreatePath(base.Config().Store.Path)
	if err != nil {
		return err
	}
	base.NodeDB = store.NewDB()
	base.NodeRaft, err = cluster.StartNode(len(base.Config().Join.Addr) == 0,
		fmt.Sprintf("%s:%v", base.Config().Server.Addr, base.Config().Server.TcpPort),
		fmt.Sprintf("%s:%v", base.Config().Server.Addr, base.Config().Server.HttpPort),
		fmt.Sprintf("%s/log", base.Config().Store.Path),
		fmt.Sprintf("%s/stable", base.Config().Store.Path),
		fmt.Sprintf("%s/snapshot", base.Config().Store.Path))
	if err != nil {
		return err
	}
	return nil
}

// JoinCluster 加入集群
func JoinCluster() error {
	if len(base.Config().Join.Addr) == 0 {
		return nil
	}
	marshal, err := json.Marshal(base.Config().Server)
	if err != nil {
		return err
	}
	_, err = base.HttpPost(fmt.Sprintf("http://%s:%s/node", base.Config().Join.Addr, base.Config().Join.HttpPort), marshal)
	if err != nil {
		return err
	}
	return nil
}

// AddNode 在领导者节点上添加新的节点
func AddNode(addr string, tcpPort int, httpPort int) error {
	if len(addr) == 0 {
		return errors.New("len(addr) == 0")
	}
	if tcpPort <= 0 {
		return errors.New("tcpPort <= 0")
	}
	if httpPort <= 0 {
		return errors.New("httpPort <= 0")
	}
	return cluster.AddNode(base.NodeRaft, fmt.Sprintf("%s:%v", addr, tcpPort), fmt.Sprintf("%s:%v", addr, httpPort))
}

func RemoveNode(endpoint string) error {
	if len(endpoint) == 0 {
		return errors.New("len(endpoint) == 0")
	}
	return cluster.RemoveNode(base.NodeRaft, endpoint)
}

func GetLeader() string {
	return cluster.GetLeader(base.NodeRaft)
}

func ListNode() []cluster.NodeInfoModel {
	return cluster.ListNode(base.NodeRaft)
}
