package core

import (
	"fraisedb/base"
	"fraisedb/cluster"
	"fraisedb/store"
	jsoniter "github.com/json-iterator/go"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func ListNamespace() []string {
	return base.NodeDB.ListNamespace()
}

func CreateNamespace(namespace string) error {
	if len(namespace) == 0 {
		return errors.New("len(namespace) == 0")
	}
	return cluster.ApplyLog(base.NodeRaft, namespace, 11, "", false, "", 0, 0)
}

func DeleteNamespace(namespace string) error {
	if len(namespace) == 0 {
		return errors.New("len(namespace) == 0")
	}
	return cluster.ApplyLog(base.NodeRaft, namespace, 10, "", false, "", 0, 0)
}

func GetKV(namespace string, key string) (*store.KvDTO, error) {
	if len(namespace) == 0 {
		return nil, errors.New("len(namespace) == 0")
	}
	if len(key) == 0 {
		return nil, errors.New("len(key) == 0")
	}
	return base.NodeDB.GetKV(namespace, key)
}

func GetNestedField(namespace string, key string, path []any) (any, error) {
	if len(namespace) == 0 {
		return nil, errors.New("len(namespace) == 0")
	}
	if len(key) == 0 {
		return nil, errors.New("len(key) == 0")
	}
	kv, err := base.NodeDB.GetKV(namespace, key)
	if err != nil {
		return nil, err
	}
	var value any
	jsoniter.Get([]byte(kv.Value), path).ToVal(&value)
	return value, nil
}

func PutKV(namespace string, key string, overwrite bool, value string, incr int64, ddl int64) error {
	if len(namespace) == 0 {
		return errors.New("len(namespace) == 0")
	}
	if len(key) == 0 {
		return errors.New("len(key) == 0")
	}
	if len(value) == 0 && overwrite {
		return errors.New("len(value) == 0")
	}
	if base.NodeDB.NamespaceNotExist(namespace) {
		return errors.New("namespace not exist")
	}
	return cluster.ApplyLog(base.NodeRaft, namespace, 1, key, overwrite, value, incr, ddl)
}

func DeleteKV(namespace string, key string) error {
	if len(namespace) == 0 {
		return errors.New("len(namespace) == 0")
	}
	if len(key) == 0 {
		return errors.New("len(key) == 0")
	}
	return cluster.ApplyLog(base.NodeRaft, namespace, 0, key, false, "", 0, 0)
}

func ListKV(namespace string, keyPrefix string, offset int64, count int64) (*[]store.KvDTO, error) {
	if len(namespace) == 0 {
		return nil, errors.New("len(namespace) == 0")
	}
	return base.NodeDB.ListKV(namespace, keyPrefix, offset, count)
}
