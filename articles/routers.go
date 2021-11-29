package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/triostones/golang-gin-realworld/common"
)

func ArticleCreate(c *gin.Context) {
	articleModelValidator := NewArticleModelValidator()
	if err := articleModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	if err := SaveOne(&articleModelValidator.articleModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ArticleSerializer{c, articleModelValidator.articleModel}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}

func ArticlesRegister(router *gin.RouterGroup) {
	router.POST("/", ArticleCreate)
	// TODO:
}

func ArticlesAnonymousRegister(router *gin.RouterGroup) {
	// TODO:
}
