// пакет для вывода html всем подписчикам
package sse

import (
	"github.com/gin-gonic/gin"

	"chicha/Packages/race"
)

func Setup(r *gin.RouterGroup, ch <-chan race.ID) {
	// Initialize new streaming server
	server := NewServer()
	server.Notifier = ch

	// client can listen stream events from server
	r.GET("/:raceID", server.serveHTTP)
}
