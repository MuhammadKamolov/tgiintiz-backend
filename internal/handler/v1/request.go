package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"
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

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func parsePagination(r *http.Request) models.Pagination {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	return models.NewPagination(page, pageSize)
}
