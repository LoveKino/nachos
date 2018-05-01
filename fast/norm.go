package fast

import (
	. "github.com/LoveKino/nachos/http"
	. "github.com/LoveKino/nachos/store"
)

const AUTH_ERROR_CODE = 50
const API_PATH_PREFIX = "/api/"
const OMIT_API_RESPONSE = "{\"errNo\": -1, \"errMsg\": \"Not Support Api\"}"

type CommonModuleConfig struct {
	SessionKey []byte
	StoreCons  StoreConstructor
}

type ApiOptions struct {
	QueryMap map[string]string
	ExecMap  map[string]string
}

type ParamsHandler func(params ...interface{}) ApiHandler

type CommonApiUtil struct {
	QuickApi        func(options ApiOptions) ParamsHandler
	QuickLoginedApi func(options ApiOptions) ParamsHandler
}
