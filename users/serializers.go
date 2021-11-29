package users

import (
	"github.com/gin-gonic/gin"
	"github.com/triostones/golang-gin-realworld/common"
)

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {

	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    myUserModel.Image,
		Token:    common.GenToken(myUserModel.ID),
	}
	return user
}

type ProfileSerializer struct {
	c *gin.Context
	UserModel
}

type ProfileResponse struct {
	ID        uint    `json:"-"`
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

func (self *ProfileSerializer) Response() ProfileResponse {
	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID: self.ID,
		Username: self.Username,
		Bio: self.Bio,
		Image: self.Image,
		Following: myUserModel.isFollowing(self.UserModel),
	}
	return profile
}
