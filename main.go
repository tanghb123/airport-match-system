package main

import (
	"airport-match-system/initial"
	"airport-match-system/log"
	"airport-match-system/router"
)

func main() {
	initial.InitConfig()
	initial.InitMysqlDB()
	initial.InitRedis()
	log.LogInit("info")
	r := router.Router()

	r.Run(":8081")
}
