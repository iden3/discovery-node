package endpoint

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/node"
)

var serverConfig config.Config
var nodesrv *node.NodeSrv

func newApiService() *gin.Engine {
	api := gin.Default()
	api.Use(cors.Default())
	api.GET("/info", handleInfo)
	api.POST("/id", handleStoreId)
	api.GET("/id/:idaddr", handleDiscoverId)
	return api
}

func Serve(cnfg config.Config, nodeservice node.NodeSrv) *gin.Engine {
	serverConfig = cnfg
	nodesrv = &nodeservice
	return newApiService()
}
