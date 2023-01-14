package main

import (
	"github.com/athunlal/config"
	"github.com/athunlal/routes"

	"github.com/athunlal/initializer"
	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnv()
	config.DBconnect()
}

var R = gin.Default()

func main() {
	gin.SetMode(gin.ReleaseMode)

	routes.AdminRouts(R)
	routes.UserRouts(R)

	R.Run()
}
