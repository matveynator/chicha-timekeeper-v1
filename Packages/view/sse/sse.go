// пакет для вывода html всем подписчикам
package sse

import (
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.RouterGroup, ch <-chan struct{}) {
	// Initialize new streaming server
	server := NewServer()
	server.Notifier = ch

	// client can listen stream events from server
	r.GET("/:id", server.serveHTTP)
}
