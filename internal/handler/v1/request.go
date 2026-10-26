package v1

import (
	"encoding/json"
	"net/http"
	"tgiintiz/internal/handler/response"
)

const maxBodyBytes = 1 << 20

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid JSON body", err)
		return false
	}

	return true
}
