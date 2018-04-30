package simple

import (
	. "../ctx"
)

/**
 * A special validation mid which used to do some assertions
 */
type MidValidator func(ctx *ApiContext, params []interface{}) error
