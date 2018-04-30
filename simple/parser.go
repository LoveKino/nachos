package simple

import (
	"strings"
)

// split into params. logic and result parts
func Parse(args []interface{}) (interface{}, interface{}, interface{}) {
	var params interface{}
	var logic interface{}
	var result interface{}

	if len(args) == 1 {
		text := args[0].(string)
		parts := strings.Split(text, "->")
		length := len(parts)
		if length >= 1 {
			params = parts[0]
		}
		if length >= 2 {
			logic = parts[1]
		}
		if length >= 3 {
			result = parts[2]
		}
	} else if len(args) == 2 {
		params = args[0]
		logic = args[1]
	} else if len(args) == 3 {
		params = args[0]
		logic = args[1]
		result = args[2]
	}

	return params, logic, result
}
