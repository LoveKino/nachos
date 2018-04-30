package http

import (
	"net/http"
)

const SERVER_ERROR_NUMBER = 1
const QUERY_ERROR_NUMBER = 100
const PARAM_MISSING_ERROR_NUMBER = 200
const POST_BODY_PARSE_ERROR_NUMBER = 300
const PARAM_ERROR_NUMBER = 300
const ERROR_ROWS_CONVERT = 500
const RESPONSE_PARSE_ERROR = 400

const QUERY_ERROR_MESSAGE = "Error happend when query from database."
const ROWS_CONVERT_ERROR_MESSAGE = "Fail to convert queried rows to target data structure"
const EMPTY_DATA = "{}"

type ApiError struct {
	ErrNo  int
	ErrMsg string
}

/**
 * export api map
 */
type ApiHandler func(http.ResponseWriter, *http.Request)
