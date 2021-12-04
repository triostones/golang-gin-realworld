package articles

import (
	"github.com/jinzhu/gorm"
	"github.com/triostones/golang-gin-realworld/common"
	"github.com/triostones/golang-gin-realworld/users"
)

type FavoriteModel struct {
	gorm.Model
	Favorite     ArticleModel
	FavoriteID   uint
	FavoriteBy   ArticleUserModel
	FavoriteByID uint
}

type ArticleUserModel struct {
	gorm.Model
	UserModel      users.UserModel
	UserModelID    uint
	ArticleModels  []ArticleModel  `gorm:"ForeignKey:AuthorID"`
	FavoriteModels []FavoriteModel `gorm:"ForeignKey:FavoriteByID"`
}

type TagModel struct {
	gorm.Model
	Tag           string         `gorm:"unique_index"`
	ArticleModels []ArticleModel `gorm:"many2many:article_tags"`
}

type CommentModel struct {
	gorm.Model
	Article   ArticleModel
	ArticleID uint
	Author    ArticleUserModel
	AuthorID  uint
	Body      string `gorm:"size:2048"`
}

type ArticleModel struct {
	gorm.Model
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	Author      ArticleUserModel
	AuthorID    uint
	Tags        []TagModel     `gorm:"many2many:article_tags"`
	Comments    []CommentModel `gorm:"ForeignKey:ArticleID"`
}

func (model *ArticleModel) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}

func (model *ArticleModel) setTags(tags []string) error {
	db := common.GetDB()
	var tagList []TagModel
	for _, tag := range tags {
		var tagModel TagModel
		err := db.FirstOrCreate(&tagModel, TagModel{Tag: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	model.Tags = tagList
	return nil
}

func (article *ArticleModel) isFavoriteBy(user ArticleUserModel) bool {
	db := common.GetDB()
	var favorite FavoriteModel
	db.Where(FavoriteModel{
		FavoriteID:   article.ID,
		FavoriteByID: user.ID,
	}).First(&favorite)
	return favorite.ID != 0
}

func (article *ArticleModel) favoritesCount() uint {
	db := common.GetDB()
	var count uint
	db.Model(&FavoriteModel{}).Where(FavoriteModel{FavoriteID: article.ID}).Count(&count)
	return count
}

func GetArticleUserModel(userModel users.UserModel) ArticleUserModel {
	var articleUserModel ArticleUserModel
	if userModel.ID == 0 {
		return articleUserModel
	}
	db := common.GetDB()
	db.Where(&ArticleUserModel{UserModelID: userModel.ID}).FirstOrCreate(&articleUserModel)
	articleUserModel.UserModel = userModel
	return articleUserModel
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func FindOneArticle(condition interface{}) (ArticleModel, error) {
	db := common.GetDB()
	var model ArticleModel
	tx := db.Begin()
	tx.Where(condition).First(&model)
	tx.Model(&model).Related(&model.Author, "Author")
	tx.Model(&model.Author).Related(&model.Author.UserModel)
	tx.Model(&model).Related(&model.Tags, "Tags")
	err := tx.Commit().Error
	return model, err
}
