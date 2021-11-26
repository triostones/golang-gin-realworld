package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
}

func main() {
	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()
	addPingGroup(r)

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))

	r.Run()
}
