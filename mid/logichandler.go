package mid

import (
	. "../ctx"
	"../util"
	"errors"
	"log"
	"reflect"
	"strings"
)

// TODO support logics, chain group of logic
func ParseLogicMidHandler(logic interface{}) AtomMidLogicHandler {
	if logic == nil {
		return EmptyLogicHandler
	}
	switch p2 := logic.(type) {
	case string:
		return PanicMidLogicHandler(p2)
	case func(*ApiContext, []interface{}, func(interface{}, error)):
		return p2
	// special case when do not need next callback
	case func(*ApiContext, []interface{}) (interface{}, error):
		return func(ctx *ApiContext, params []interface{}, next func(interface{}, error)) {
			ret, err := p2(ctx, params)
			next(ret, err)
		}
	default:
		log.Print(reflect.TypeOf(p2))
		panic(errors.New("type of logic for SimpleAtomMid must be string or AtomMidLogicHandler."))
	}
}

/**
 * quick logic handler
 */
func PanicMidLogicHandler(action string) AtomMidLogicHandler {
	handler, err := QuickMidLogicHandler(action)
	if err != nil {
		panic(err)
	} else {
		return handler
	}
}

func QuickMidLogicHandler(action string) (AtomMidLogicHandler, error) {
	parts := strings.Split(action, ":")
	actionType := strings.TrimSpace(parts[0])
	actionCnt := ""
	if len(parts) > 1 {
		actionCnt = strings.TrimSpace(parts[1])
	}

	if actionType == "query" { // query store
		// TODO validate query first
		return func(ctx *ApiContext, params []interface{}, next func(interface{}, error)) {
			m, qerr := ctx.Store.Querys[actionCnt](params...)
			next(m, qerr)
		}, nil
	} else if actionType == "exec" { // exec store
		return func(ctx *ApiContext, params []interface{}, next func(interface{}, error)) {
			ret, exeErr := ctx.Store.Execs[actionCnt](params...)
			next(ret, exeErr)
		}, nil
	} else if actionType == "pass" { // pass param
		return func(ctx *ApiContext, params []interface{}, next func(interface{}, error)) {
			text := strings.TrimSpace(actionCnt)
			if text == "" {
				text = "0"
			}
			v, geterr := util.GetValueByJsonPath(params, text)
			next(v, geterr)
		}, nil
	}
	// TODO more common logic types

	return nil, errors.New("there is no logic handler type " + actionType)
}

func EmptyLogicHandler(ctx *ApiContext, params []interface{}, next func(interface{}, error)) {
	// do nothing, pass power to next
	next(nil, nil)
}
