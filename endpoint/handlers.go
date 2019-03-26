package endpoint

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/iden3/discovery-node/discovery"
	log "github.com/sirupsen/logrus"
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
	color.Cyan("id stored: " + id.IdAddr.Hex())
	if err != nil {
		fail(c, "error storing id", err)
	}

	c.JSON(200, gin.H{})
}

func handleDiscoverId(c *gin.Context) {
	idAddrStr := c.Param("idaddr")
	idAddr := common.HexToAddress(idAddrStr)
	_, err := nodesrv.DiscoverId(idAddr)
	if err != nil {
		fail(c, "error storing id", err)
	}
	// color.Cyan("id discovered: " + id.IdAddr.Hex())
	c.JSON(200, gin.H{
		"foo": "bar",
	})
}
