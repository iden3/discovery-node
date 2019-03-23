package endpoint

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/iden3/discovery-research/discovery-node/discovery"
)

func fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

func handleInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "bar",
	})
}

func handleStoreId(c *gin.Context) {
	var id discovery.Id
	c.BindJSON(&id)


	err := nodesrv.StoreId(id)
	if err!=nil {
		fail(c, "error storing id", err)
	}
	
	c.JSON(200, gin.H{})
}

func handleDiscoverId(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "bar",
	})
}
