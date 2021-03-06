package fast

import (
	. "github.com/LoveKino/nachos/ctx"
	. "github.com/LoveKino/nachos/http"
	. "github.com/LoveKino/nachos/simple"
	. "github.com/LoveKino/nachos/store"
	"log"
	"net/http"
	"strings"
)

// generate api handler by ApiHandlerMap
func MapToApiHandler(apiMap ApiHandlerMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		log.Printf("Entrance: " + urlPath + "?" + r.URL.RawQuery)

		if strings.HasPrefix(urlPath, API_PATH_PREFIX) {
			// pick up api function
			apiHandler, hasApi := apiMap[urlPath]

			if hasApi {
				log.Printf("Run api: " + urlPath)
				apiHandler(w, r)
			} else {
				w.Write([]byte(OMIT_API_RESPONSE))
			}
		} else {
			w.Write([]byte(OMIT_API_RESPONSE))
		}
	}
}

var LoginCheck = SimpleAtomMid("session:sid", func(ctx *ApiContext, params []interface{}) (interface{}, error) {
	if sid := params[0].(string); sid != "" {
		return sid, nil
	} else {
		return nil, ApiError{AUTH_ERROR_CODE, "login failure, please relogin"}
	}
}, "setMidMap:userId")

func CommonHttpModule(moduleConfig CommonModuleConfig) CommonApiUtil {
	quickApi := func(options ApiOptions) ParamsHandler {
		return func(params ...interface{}) ApiHandler {
			store := GetStore(moduleConfig.StoreCons, options.QueryMap, options.ExecMap)
			return ToHttpHandler(store, ApiConfig{
				"SessionKey": moduleConfig.SessionKey,
			}, SeqAtomMids(params...))
		}
	}

	quickLoginedApi := func(options ApiOptions) ParamsHandler {
		qApi := quickApi(options)
		return func(params ...interface{}) ApiHandler {
			// append login check
			nextParams := append([]interface{}{LoginCheck}, params...)
			return qApi(nextParams...)
		}
	}

	apiUtil := CommonApiUtil{
		QuickApi:        quickApi,
		QuickLoginedApi: quickLoginedApi,
	}

	return apiUtil
}
