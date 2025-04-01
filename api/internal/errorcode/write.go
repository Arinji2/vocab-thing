package errorcode

import (
	"encoding/json"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	if appErr, ok := err.(*AppError); ok {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]any{
			"errorCode": appErr.Code,
			"message":   appErr.Message,
			"readable":  appErr.Readable,
		})
		return
	}

	// Handle unexpected errors
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]any{
		"errorCode": 500,
		"message":   "Internal Server Error",
		"readable":  "Internal Server Error",
		"details":   err.Error(),
	})
}
