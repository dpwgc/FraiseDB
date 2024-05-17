package core

import (
	"fmt"
	"fraisedb/base"
	"fraisedb/cluster"
	"fraisedb/store"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// InitNode 启动节点
func InitNode() error {
	base.Channel = make(chan []byte, 1000)
	err := base.CreatePath(base.Config().Store.Data)
	if err != nil {
		return err
	}
	base.NodeDB = store.NewDB()
	base.NodeRaft, err = cluster.StartNode(base.Config().Node.First,
		fmt.Sprintf("%s:%v", base.Config().Node.Addr, base.Config().Node.TcpPort),
		fmt.Sprintf("%s:%v", base.Config().Node.Addr, base.Config().Node.HttpPort),
		fmt.Sprintf("%s/log", base.Config().Store.Data),
		fmt.Sprintf("%s/stable", base.Config().Store.Data),
		fmt.Sprintf("%s/snapshot", base.Config().Store.Data))
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
