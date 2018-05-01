package simple

import (
	. "github.com/LoveKino/nachos/ctx"
	. "github.com/LoveKino/nachos/mid"
)

// TODO name mapping for row maps
// TODO type mapping for row maps

/**
 * a simple generic way to define mids
 */
func SimpleAtomMid(args ...interface{}) AtomMid {
	params, logic, result := Parse(args)

	return ToAtomMid(ParseParamsCollector(params),
		ParseLogicMidHandler(logic),
		ParseResultHandler(result))
}

/**
 * validate some params, if error, handle error (default handle is flush error)
 *          if pass, do nothing
 */
func ValidateMid(params interface{}, validator MidValidator) AtomMid {
	return SimpleAtomMid(params, func(ctx *ApiContext, args []interface{}) (interface{}, error) {
		err := validator(ctx, args)
		return nil, err
	}, "flush:errorOnly")
}

// just flush
func JustFlush(result interface{}) AtomMid {
	return SimpleAtomMid([]string{}, nil, result)
}
