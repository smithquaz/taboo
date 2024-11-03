package routes

import (
	"taboo-game/models"

	"github.com/gin-gonic/gin"
)

type WordRoutes struct {
	commonWords []models.Word
	domainWords []models.Word
}

func NewWordRoutes(commonWords, domainWords []models.Word) *WordRoutes {
	return &WordRoutes{
		commonWords: commonWords,
		domainWords: domainWords,
	}
}

func (wr *WordRoutes) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/words", wr.getWords)
}

func (wr *WordRoutes) getWords(c *gin.Context) {
	category := c.Query("category")
	switch category {
	case "common":
		c.JSON(200, wr.commonWords)
	case "domain":
		c.JSON(200, wr.domainWords)
	default:
		c.JSON(200, gin.H{
			"common":   wr.commonWords,
			"specific": wr.domainWords,
		})
	}
}
