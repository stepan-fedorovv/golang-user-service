package http_utils

import (
	"encoding/json"
	"net/http"
)

func EncodeBody(w http.ResponseWriter, s interface{}) error {
	return json.NewEncoder(w).Encode(s)
}
