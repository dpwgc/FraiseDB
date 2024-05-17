package cluster

import (
	"bytes"
	"encoding/gob"
	"fraisedb/base"
	"github.com/hashicorp/raft"
)

type Snapshot struct {
}

func newSnapshot() raft.FSMSnapshot {
	return &Snapshot{}
}

type KVSnapshotModel struct {
	Namespace string
	Key       string
	Value     string
	DDL       int64
}

// Persist saves the FSM snapshot out to the given sink.
func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := error(nil)
	defer func() {
		if err != nil {
			base.LogHandler.Println(base.LogErrorTag, err)
			err = sink.Cancel()
			if err != nil {
				base.LogHandler.Println(base.LogErrorTag, err)
			}
		}
	}()
	var kvSnaps []KVSnapshotModel
	ns := base.NodeDB.ListNamespace()
	for _, n := range ns {
		kvs, err := base.NodeDB.ListKV(n, "", 0, 0)
		if err != nil {
			return err
		}
		for _, kv := range *kvs {
			kvSnap := KVSnapshotModel{
				Namespace: n,
				Key:       kv.Key,
				Value:     kv.Value,
				DDL:       kv.DDL,
			}
			kvSnaps = append(kvSnaps, kvSnap)
		}
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(kvSnaps)
	if err != nil {
		return err
	}
	if _, err = sink.Write(buffer.Bytes()); err != nil {
		return err
	}
	if err = sink.Close(); err != nil {
		return err
	}
	return nil
}
func (s *Snapshot) Release() {}
