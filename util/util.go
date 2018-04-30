package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// get value by json path
// type explaination
// array:  []interface{}
// object: map[string]interface{}
// atom:   interface{}
func GetValueByJsonPath(value interface{}, jsonPath string) (interface{}, error) {
	var currentObject = value

	parts := strings.Split(jsonPath, ".")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			// check if it is number
			num, toNumErr := strconv.Atoi(part)
			if toNumErr == nil {
				// try array
				nextObjectParent, typeOk := currentObject.([]interface{})
				if typeOk {
					if num < 0 || num > len(nextObjectParent) {
						return nil, errors.New("missing value for path: " + jsonPath + ". Out of range. Array length is " + strconv.Itoa(len(nextObjectParent)) + ".")
					}
					nextObject := nextObjectParent[num]
					currentObject = nextObject
					continue
				}
			}

			// otherwise regarding as map
			nextObjectParent, typeOk := currentObject.(map[string]interface{})
			if !typeOk {
				return nil, errors.New("Can not go deeper for this jsonPath: " + jsonPath + ". Type of current object is " + fmt.Sprintf("%v", reflect.TypeOf(currentObject)))
			}

			nextObject, getOk := nextObjectParent[part]

			if !getOk {
				return nil, errors.New("missing value for path: " + jsonPath)
			} else {
				currentObject = nextObject
			}
		}
	}

	return currentObject, nil
}
