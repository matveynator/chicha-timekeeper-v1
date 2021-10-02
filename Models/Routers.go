package Models

/**
* Module have methods for internal routing based on Gin framework
 */

import (
	"chicha/Packages/view"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api")
	{
		v1.GET("laps", GetListLaps) // Full list of laps
		v1.GET("laps/:id", GetLap)  // Get one lap details
		v1.GET("laps/bytagid/:id", GetLapsByTagId)
		v1.GET("laps/byraceid/:id", GetLapsByRaceId)
		v1.GET("laps/results/byraceid/:id", GetResultsByRaceId)
		//v1.GET("laps/delete/bytagid/:id", DeleteLap)
		v1.GET("laps/last", GetLastLapData)

		//v1.GET("races", GetListRaces)
		//v1.GET("races/:id", GetRace)
		//v1.POST("races", CreateRace)
		//v1.PUT("races/:id", UpdateRace)
		//v1.DELETE("races/:id", DeleteRace)
		//v1.POST("races/start/:id", StartRace)
		//v1.POST("races/finish", FinishRace)

		//v1.GET("users", GetListUsers)
		//v1.GET("users/:id", GetUser)
		//v1.POST("users", CreateUser)
		//v1.PUT("users/:id", UpdateUser)
		//v1.DELETE("users/:id", DeleteUser)

		//v1.GET("checkins", GetListCheckins)
		//v1.GET("checkins/:id", GetCheckin)
		//v1.POST("checkins", CreateCheckin)
		//v1.PUT("checkins/:id", UpdateCheckin)
		//v1.DELETE("checkins/:id", DeleteCheckin)

		//v1.GET("admins", GetListAdmins)
		//v1.GET("admins/:id", GetAdmin)
		//v1.POST("admins", CreateAdmin)
		//v1.PUT("admins/:id", UpdateAdmin)
		//v1.DELETE("admins/:id", DeleteAdmin)
		//v1.POST("admins/login", LoginAdmin)

	}

	v := view.New()
	{
		r.GET("/", v.Homepage)
	}

	return r
}
