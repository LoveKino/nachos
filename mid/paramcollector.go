package mid

import (
	"errors"
	"github.com/satori/go.uuid"
	. "nachos/ctx"
	. "nachos/http"
	"nachos/session"
	"nachos/util"
	"strconv"
	"strings"
)

func ParseParamsCollector(params interface{}) AtomMidParamsCollector {
	if params == nil {
		return EmptyParamsCollector
	}

	switch p1 := params.(type) {
	case string:
		return PanicFetchParams(strings.Split(p1, ";"))
	case []string:
		return PanicFetchParams(p1)
	case func(*ApiContext) ([]interface{}, error):
		return p1
	default:
		panic(errors.New("type of params for SimpleAtomMid must be string or []string or AtomMidParamsCollector."))
	}
}

/**
 * define collectors which used to fetch params from ctx
 */
type SingleCollector func(*ApiContext) (interface{}, error)

func PanicFetchParams(addrs []string) AtomMidParamsCollector {
	collector, err := FetchCtxParams(addrs)
	if err != nil {
		panic(err)
	} else {
		return collector
	}
}

func FetchCtxParams(addrs []string) (AtomMidParamsCollector, error) {
	var collectors []SingleCollector
	for _, addr := range addrs {
		addr = strings.TrimSpace(addr)
		if addr != "" { // ignore empty string
			collector, err := FetchCtxParam(addr)
			if err != nil {
				return nil, err
			} else {
				collectors = append(collectors, collector)
			}
		}
	}

	return func(ctx *ApiContext) ([]interface{}, error) {
		params := []interface{}{}
		for _, collector := range collectors {
			param, err := collector(ctx)
			if err != nil {
				return nil, err
			} else {
				params = append(params, param)
			}
		}
		return params, nil
	}, nil
}

/**
 * help function which used to get param from context
 * addr = type : address
 * TODO support expandation?
 */
func FetchCtxParam(addr string) (func(*ApiContext) (interface{}, error), error) {
	parts := strings.Split(addr, ":")
	addrType := strings.TrimSpace(parts[0])
	addrCnt := ""
	if len(parts) > 1 {
		addrCnt = strings.TrimSpace(parts[1])
	}

	if addrType == "url" { // get param from url
		return func(ctx *ApiContext) (interface{}, error) {
			value, ok := GetUrlParam(ctx.Request, addrCnt)
			if !ok {
				return nil, errors.New("missing " + addrCnt + " in url param")
			} else {
				return value, nil
			}
		}, nil
	} else if addrType == "post" { // get post data attribute
		return func(ctx *ApiContext) (interface{}, error) {
			duck, err := getCachePostBody(ctx)
			if err != nil {
				return nil, err
			} else {
				value, err := util.GetValueByJsonPath(duck, addrCnt)
				if err != nil {
					return nil, errors.New("missing " + addrCnt + " in post body. Error details: " + err.Error())
				} else {
					return value, nil
				}
			}
		}, nil
	} else if addrType == "cookie" { // get cookie value
		return func(ctx *ApiContext) (interface{}, error) {
			cookie, err := ctx.Request.Cookie(addrCnt)
			if err != nil {
				return "", nil
			}
			return cookie.Value, nil
		}, nil
	} else if addrType == "session" { // get session value
		return func(ctx *ApiContext) (interface{}, error) {
			sessionKey := ctx.Config["SessionKey"].([]byte)
			return session.GetSession(ctx.Request, sessionKey, addrCnt), nil
		}, nil
	} else if addrType == "midMap" { // get value from midMap
		return func(ctx *ApiContext) (interface{}, error) {
			// get param from mid
			value, err := util.GetValueByJsonPath(ctx.MidMap, addrCnt)
			if err != nil {
				return nil, errors.New("missing " + addrCnt + " in midmap. Error details: " + err.Error())
			} else {
				return value, nil
			}
		}, nil
	} else if addrType == "const_str" { // string constant
		// TODO encoding string content
		return func(ctx *ApiContext) (interface{}, error) {
			return addrCnt, nil
		}, nil
	} else if addrType == "const_int" { // int constant
		intV, siErr := strconv.Atoi(addrCnt)
		if siErr != nil {
			return nil, siErr
		} else {
			return func(ctx *ApiContext) (interface{}, error) {
				return intV, nil
			}, nil
		}
		// TODO more constant type
	} else if addrType == "uuid" {
		return func(ctx *ApiContext) (interface{}, error) {
			return uuid.NewV4()
		}, nil
	}
	// TODO more types

	return nil, errors.New("there is no param fetcher type: " + addrType)
}

func EmptyParamsCollector(ctx *ApiContext) ([]interface{}, error) {
	return nil, nil
}

func getCachePostBody(ctx *ApiContext) (map[string]interface{}, error) {
	if ctx.PostBodyCache == nil {
		duck, err := JsonBodyToMap(ctx.Request)
		if err != nil {
			return nil, err
		}

		ctx.PostBodyCache = duck
	}

	return ctx.PostBodyCache, nil
}
