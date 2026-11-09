package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"tgiintiz/internal/handler/response"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				err := fmt.Errorf("panic: %v\n%s", p, debug.Stack())
				response.Error(c, http.StatusInternalServerError, "internal server error", err)
			}
		}()
		c.Next()
	}
}
