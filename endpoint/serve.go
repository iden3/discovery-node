package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/iden3/discovery-research/discovery-node/config"
)

var serverConfig config.Config

func newApiService() *gin.Engine {
	api := gin.Default()
	api.Use(cors.Default())
	api.GET("/info", handleInfo)
	api.POST("/id", handleStoreId)
	api.GET("/id/:idaddr", handleDiscoverId)
	return api
}

func Serve(cnfg config.Config) *gin.Engine {
	serverConfig = cnfg
	return newApiService()
}
