package ctx

import (
	. "nachos/http"
	. "nachos/store"
	"net/http"
)

/**
 * standard API signature
 *
 * (ctx *ApiContext) -> Void
 */

type ApiContext struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Store         Store
	PostBodyCache map[string]interface{}
	MidMap        map[string]interface{} // provide for mids usage

	// TODO use map
	Config ApiConfig
}

type ApiConfig map[string]interface{}

type ApiHandlerMap map[string]ApiHandler

func ToHttpHandler(store Store, config ApiConfig, api func(*ApiContext)) ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		MidMap := make(map[string]interface{})
		api(&ApiContext{
			w,
			r,
			store,
			nil,
			MidMap,
			config,
		})
	}
}
