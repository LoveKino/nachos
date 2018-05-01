package simple

import (
	. "github.com/LoveKino/nachos/ctx"
)

/**
 * A special validation mid which used to do some assertions
 */
type MidValidator func(ctx *ApiContext, params []interface{}) error
