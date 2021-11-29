package users

import (
	"errors"
	"github.com/triostones/golang-gin-realworld/common"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint    `gorm:primary_key`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

func (u *UserModel) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil

}

func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashPassword, bytePassword)
}

func (model *UserModel) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}

func AutoMigrate() {
	db := common.Init()
	db.AutoMigrate(&UserModel{})
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}
