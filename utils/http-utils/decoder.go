package http_utils

import (
	"encoding/json"
	"net/http"
)

func DecodeBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(s)
	return err
}
