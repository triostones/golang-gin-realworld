package users

import (
	"github.com/gin-gonic/gin"
	"github.com/triostones/golang-gin-realworld/common"
)

type UserModelValidator struct {
	User struct {
		Username string `form:"username" json:"username" binding:"required,alphanum,min=4,max=25"`
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required,min=8,max=255"`
		Bio      string `form:"bio" json:"bio" binding:"max=1024"`
		Image    string `form:"image" json:"image" binding:"omitempty,url"`
	} `json:"user"`
	userModel UserModel `json:"-"`
}

func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}

func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}

	self.userModel.Username = self.User.Username
	self.userModel.Email = self.User.Email
	self.userModel.Bio = self.User.Bio
	if self.User.Image != "" {
		self.userModel.Image = &self.User.Image
	}
	if self.User.Password != common.NBRandomPassword {
		self.userModel.setPassword(self.User.Password)
	}
	return nil
}
