package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
)

func shortenString(origin []byte, maxLen int) string {
	shortened := string(origin)
	if len(shortened) > maxLen {
		shortened = shortened[:maxLen]
	}
	return shortened
}
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *models.User
		u, exist := c.Get("user")
		if !exist {
			user = &models.User{
				BaseModel: models.BaseModel{ID: 0},
				Email:     "not-set-yet",
			}
		} else {
			user = u.(*models.User)
		}

		// 请求
		reqBody, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // 把读过的body再放回去
		fmt.Printf("\n[%s] (%v, %v) Request: %s %s\nHeaders: %v\nBody: %s\n",
			time.Now().Format(time.RFC3339),
			user.ID, user.Email,
			c.Request.Method,
			c.Request.URL,
			c.Request.Header,
			shortenString(reqBody, 1000),
		)

		// 响应
		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		responseBody := writer.body.String()
		fmt.Printf("\n[%s] (%v, %v) (%v, %v) Response: %d %s\nHeaders: %v\nBody: %s\n",
			time.Now().Format(time.RFC3339),
			user.ID, user.Email,
			c.Request.Method, c.Request.URL,
			writer.Status(),
			http.StatusText(writer.Status()),
			writer.Header(),
			responseBody,
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
