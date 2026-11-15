package v1

import (
	"net/http"
	"strconv"
	"tgiintiz/internal/handler/response"
	"tgiintiz/internal/models"

	"github.com/gin-gonic/gin"
)

func decodeJSON(c *gin.Context, dst any) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	if err := c.ShouldBindJSON(dst); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid JSON body", err)
		return false
	}

	return true
}

func parseID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func parsePagination(c *gin.Context) models.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	return models.NewPagination(page, pageSize)
}
