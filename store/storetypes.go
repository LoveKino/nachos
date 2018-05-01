package store

type Records []interface{}

type QueryStoreType func(params ...interface{}) (Records, error)
type ExecStoreType func(params ...interface{}) (interface{}, error)

type QueryStoreMap map[string]QueryStoreType
type ExecStoreMap map[string]ExecStoreType

type Store struct {
	Querys QueryStoreMap
	Execs  ExecStoreMap
}

type StoreConstructor struct {
	GetQueryStore func(query string) QueryStoreType
	GetExecStore  func(exec string) ExecStoreType
}
