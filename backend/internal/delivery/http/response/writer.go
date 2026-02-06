package response

import (
	"encoding/json"
	"net/http"
)

func write(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"status":false,"message":"error encoding response"}`, http.StatusInternalServerError)
	}
}
