package store

func GetStore(storeConstructor StoreConstructor, queryMap map[string]string, execMap map[string]string) Store {
	queryStoreMap := QueryStoreMap{}
	execStoreMap := ExecStoreMap{}

	for queryName, querySql := range queryMap {
		queryStoreMap[queryName] = storeConstructor.GetQueryStore(querySql)
	}

	for execName, execSqlText := range execMap {
		execStoreMap[execName] = storeConstructor.GetExecStore(execSqlText)
	}

	return Store{queryStoreMap, execStoreMap}
}
