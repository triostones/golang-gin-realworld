package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/triostones/golang-gin-realworld/common"
)

func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	if err := SaveOne(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, serializer.Response())
}

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
}
