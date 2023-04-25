package middlewares

import (
	"bytes"
	"fmt"
	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		// 请求
		reqBody, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // 把读过的body再放回去
		fmt.Printf("\n[%s] (%v, %v) Request: %s %s\nHeaders: %v\nBody: %s\n",
			time.Now().Format(time.RFC3339),
			user.ID, user.Email,
			c.Request.Method,
			c.Request.URL,
			c.Request.Header,
			string(reqBody),
		)

		// 响应
		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		responseBody := writer.body.String()
		fmt.Printf("\n[%s] (%v, %v) Response: %d %s\nHeaders: %v\nBody: %s\n",
			time.Now().Format(time.RFC3339),
			user.ID, user.Email,
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
