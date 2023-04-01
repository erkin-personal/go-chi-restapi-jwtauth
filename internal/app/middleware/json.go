
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error, errorCode string) {
	if err != nil {
		JSON(w, statusCode, struct {
			ErrorCode string `json:"resultCode"`
			Error     string `json:"result"`
		}{
			Error:     err.Error(),
			ErrorCode: errorCode,
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
