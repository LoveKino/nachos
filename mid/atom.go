package mid

import (
	. "github.com/LoveKino/nachos/ctx"
)

// provide basic concepts about mid

type AtomMidParamsCollector func(*ApiContext) ([]interface{}, error)

// (params, next)
type AtomMidLogicHandler func(*ApiContext, []interface{}, func(interface{}, error))

// handler result
type AtomMidResultHandler func(*ApiContext, interface{}, error)

type AtomMid func(ctx *ApiContext, next func())

// compose param collector, logic handler and result handler to a atom mid
func ToAtomMid(collector AtomMidParamsCollector,
	logicHandler AtomMidLogicHandler,
	resultHandler AtomMidResultHandler) func(*ApiContext, func()) {
	return func(ctx *ApiContext, next func()) {
		params, err := collector(ctx)
		if err != nil {
			// when errored, will not continue to next
			resultHandler(ctx, nil, err)
		} else {
			logicHandler(ctx, params, func(result interface{}, err error) {
				if err != nil {
					resultHandler(ctx, nil, err)
				} else {
					resultHandler(ctx, result, nil)
					next()
				}
			})
		}
	}
}
