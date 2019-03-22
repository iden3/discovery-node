package endpoint

import "github.com/gin-gonic/gin"

func handleInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "bar",
	})
}
