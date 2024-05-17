package store

func NewDB() DB {
	return newLevelDB()
}

type ValueModel struct {
	Value string
	DDL   int64
}

type KvDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	DDL   int64  `json:"ddl"`
}

type DB interface {
	ListNamespace() []string
	NamespaceNotExist(namespace string) bool
	CreateNamespace(namespace string) error
	DeleteNamespace(namespace string) error
	GetKV(namespace string, key string) (*KvDTO, error)
	PutKV(namespace string, key string, overwrite bool, value string, incr int64, ddl int64) error
	DeleteKV(namespace string, key string) error
	ListKV(namespace string, keyPrefix string, offset int64, count int64) (*[]KvDTO, error)
}
