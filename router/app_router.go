package router

import (
	"airport-match-system/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	passengerGroup := r.Group("/v1/passenger")
	{
		passengerGroup.POST("/create_route", service.PassengerCreateRoute)
		passengerGroup.POST("/find_match", service.FindMatchOrder)
		passengerGroup.POST("/execute_match", service.ExecuteMatchOrder)
	}

	driverGroup := r.Group("/v1/driver")
	{
		driverGroup.POST("/create_route", service.DriverCreateRoute)
	}

	return r
}
