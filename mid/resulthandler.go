package mid

import (
	"errors"
	. "github.com/LoveKino/nachos/ctx"
	. "github.com/LoveKino/nachos/http"
	"github.com/LoveKino/nachos/session"
	"strings"
)

func ParseResultHandler(v interface{}) AtomMidResultHandler {
	if v == nil {
		return EmptyResultHandler
	}

	var resultHandler AtomMidResultHandler
	switch p := v.(type) {
	case string:
		resultHandler = PanicMidRetHandler(p, ErrorFlush)
	case func(*ApiContext, interface{}, error):
		resultHandler = p
	default:
		panic(errors.New("type of result for SimpleAtomMid must be string or AtomMidResultHandler."))
	}

	return resultHandler
}

/**
 * some quick way to define result handler
 */

func PanicMidRetHandler(action string, errHandler func(ctx *ApiContext, err error)) AtomMidResultHandler {
	actions := strings.Split(action, ";")
	handle, err := QuickMidRetHandlers(actions, errHandler)
	if err != nil {
		panic(err)
	} else {
		return handle
	}
}

func QuickMidRetHandlers(actions []string, errHandler func(ctx *ApiContext, err error)) (AtomMidResultHandler, error) {
	var handlers []AtomMidResultHandler
	for _, action := range actions {
		action = strings.TrimSpace(action)
		if action != "" {
			handler, err := QuickMidRetHandler(action, errHandler)
			if err != nil {
				return nil, err
			} else {
				handlers = append(handlers, handler)
			}
		}
	}

	return MergeAtomMidResultHandlers(handlers), nil
}

// TODO error for atom mid result handler
func MergeAtomMidResultHandlers(handlers []AtomMidResultHandler) AtomMidResultHandler {
	return func(ctx *ApiContext, result interface{}, err error) {
		for _, handler := range handlers {
			handler(ctx, result, err)
		}
	}
}

/**
 * generate quick mid result handler by string
 * eg: "flush", "setMidMap:userId"
 */
func QuickMidRetHandler(action string, errHandler func(ctx *ApiContext, err error)) (AtomMidResultHandler, error) {
	parts := strings.Split(action, ":")
	actionType := strings.TrimSpace(parts[0])
	actionCnt := ""
	if len(parts) > 1 {
		actionCnt = strings.TrimSpace(parts[1])
	}

	if actionType == "flush" { // flush result
		if actionCnt != "" &&
			actionCnt != "empty" &&
			actionCnt != "errorOnly" {
			return nil, errors.New("unexpected content for flush, content " + actionCnt)
		}
		return func(ctx *ApiContext, data interface{}, err error) {
			if err != nil {
				errHandler(ctx, err)
			} else {
				if actionCnt == "" {
					WriteResponseData(ctx.Writer, data)
				} else if actionCnt == "empty" {
					WriteEmptyResponseData(ctx.Writer)
				} else if actionCnt == "errorOnly" {
					// Do nothing
				}
			}
		}, nil
	} else if actionType == "setMidMap" { // "setMidMap: key"
		return func(ctx *ApiContext, result interface{}, err error) {
			if err != nil {
				errHandler(ctx, err)
			} else {
				ctx.MidMap[actionCnt] = result
			}
		}, nil
	} else if actionType == "setSession" {
		return func(ctx *ApiContext, result interface{}, err error) {
			if err != nil {
				errHandler(ctx, err)
			} else {
				sessionValue := result.(string)
				fail := session.QuickSetSession(ctx.Writer, ctx.Config["SessionKey"].([]byte), actionCnt, sessionValue, 30*24*60)
				if fail != nil {
					errHandler(ctx, fail)
				}
			}
		}, nil
	}
	// TODO more types

	return nil, errors.New("there is no result handler type " + actionType)
}

func ErrorFlush(ctx *ApiContext, err error) {
	if v, ok := err.(ApiError); ok {
		WriteError(ctx.Writer, v.ErrNo, err.Error())
	} else {
		WriteError(ctx.Writer, SERVER_ERROR_NUMBER, err.Error())
	}
}

func EmptyResultHandler(*ApiContext, interface{}, error) {
	// do nothing
}
