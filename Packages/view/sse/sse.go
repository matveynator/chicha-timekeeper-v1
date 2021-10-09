package sse

import (
	"io"

	"github.com/gin-gonic/gin"
)

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func Setup(r *gin.RouterGroup, ch <-chan struct{}) {
	// Add event-streaming headers
	r.Use(HeadersMiddleware())

	// Initialize new streaming server
	stream := NewServer()
	r.Use(stream.serveHTTP())

	// client can listen stream events from server
	r.GET("/:id", streamHandler(ch))
}

func streamHandler(ch <-chan struct{}) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-ch; ok {
				c.SSEvent("update", msg)
				return true
			}
			return false
		})
	}
}
