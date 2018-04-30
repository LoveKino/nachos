package http

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (e ApiError) Error() string {
	return e.ErrMsg
}

func ToResponseData(data string) []byte {
	return []byte("{\"errNo\": 0, \"errMsg\": \"\", \"data\": " + data + "}")
}

func ToErrorResponseData(errNo int, errMsg string) []byte {
	return []byte("{\"errNo\": " + strconv.Itoa(errNo) + ", \"errMsg\": \"" + errMsg + "\"}")
}

func GetUrlParam(r *http.Request, key string) (string, bool) {
	queryMap := r.URL.Query()
	items, ok := queryMap[key]
	if ok {
		return items[0], true
	} else {
		return "", false
	}
}

func DecodePostBody(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dest)
	defer r.Body.Close()

	return err
}

/**
 * quick write apis
 */
func WriteResponseData(w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)

	if err != nil {
		w.Write(ToErrorResponseData(RESPONSE_PARSE_ERROR, err.Error()))
	} else {
		w.Write(ToResponseData(string(bytes)))
	}
}

func WriteError(w http.ResponseWriter, errNo int, errMsg string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)

	w.Write(ToErrorResponseData(errNo, errMsg))
}

// specific
func WriteMissingParamError(w http.ResponseWriter, paramName string) {
	WriteError(w, PARAM_MISSING_ERROR_NUMBER, "Missing param "+paramName)
}

func WriteQueryError(w http.ResponseWriter) {
	WriteError(w, QUERY_ERROR_NUMBER, QUERY_ERROR_MESSAGE)
}

func WritePostParseError(w http.ResponseWriter, errMsg string) {
	WriteError(w, POST_BODY_PARSE_ERROR_NUMBER, "Fail to parse post body, error infor is "+errMsg)
}

func WriteEmptyResponseData(w http.ResponseWriter) {
	WriteResponseData(w, ToResponseData(EMPTY_DATA))
}

func WriteRowsConvertError(w http.ResponseWriter) {
	WriteError(w, ERROR_ROWS_CONVERT, ROWS_CONVERT_ERROR_MESSAGE)
}
