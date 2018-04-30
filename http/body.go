package http

import (
	"encoding/json"
	"net/http"
)

func JsonBodyToMap(r *http.Request) (map[string]interface{}, error) {
	// TODO
	duck := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&duck)
	if err != nil {
		return nil, err
	} else {
		defer r.Body.Close()
		return duck, nil
	}
}
