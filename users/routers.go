package users

import (
	"errors"
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

func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := FindOneUser(&UserModel{Email: loginValidator.userModel.Email})
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("Not Registered email or invalid password"))
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusNotFound, errors.New("Not Registered email or invalid password"))
		return
	}
	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserRetrive(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserUpdate(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(UserModel)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
	}
	UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func ProfileRetrive(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)
	err = myUserModel.following(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)
	err = myUserModel.unFollowing(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrive)
	router.PUT("/", UserUpdate)
}

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrive)
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}
