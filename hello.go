package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/triostones/golang-gin-realworld/articles"
	"github.com/triostones/golang-gin-realworld/common"
	"github.com/triostones/golang-gin-realworld/users"
)

func addPingGroup(r *gin.Engine) {
	pingGroup := r.Group("/api/ping")
	pingGroup.GET(
		"/",
		func(c *gin.Context) {
			c.JSON(
				200,
				gin.H{
					"message": "pong"})
		})
}

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})
}

func main() {
	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()
	addPingGroup(r)

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	// TODO: tags

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))
	articles.ArticlesRegister(v1.Group("/articles"))

	r.Run()
}
