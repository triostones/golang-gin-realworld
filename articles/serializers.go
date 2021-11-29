package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/triostones/golang-gin-realworld/users"
)

type ArticleUserSerializer struct {
	C *gin.Context
	ArticleUserModel
}

func (s *ArticleUserSerializer) Response() users.ProfileResponse {
	response := users.ProfileSerializer{s.C, s.ArticleUserModel.UserModel}
	return response.Response()
}

type TagSerializer struct {
	C *gin.Context
	TagModel
}
func (s *TagSerializer) Response() string {
	return s.TagModel.Tag
}

type ArticleSerializer struct {
	c *gin.Context
	ArticleModel
}

type ArticleResponse struct {
	ID             uint                  `json:"-"`
	Title          string                `json:"title"`
	Slug           string                `json:"slug"`
	Description    string                `json:"description"`
	Body           string                `json:"body"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
	Author         users.ProfileResponse `json:"author"`
	Tags           []string              `json:"tagList"`
	Favorite       bool                  `json:"favorited"`
	FavoritesCount uint                  `json:"favoritesCount"`
}

func (s *ArticleSerializer) Response() ArticleResponse {
	myUserModel := s.c.MustGet("my_user_model").(users.UserModel)
	authorSerializer := ArticleUserSerializer{s.c, s.Author}
	response := ArticleResponse{
		ID:             s.ID,
		Slug:           slug.Make(s.Title),
		Title:          s.Title,
		Description:    s.Description,
		Body:           s.Body,
		CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt:      s.UpdatedAt.Format("2006-01-02T15:04:05.999Z"),
		Author:         authorSerializer.Response(),
		Favorite:       s.isFavoriteBy(GetArticleUserModel(myUserModel)),
		FavoritesCount: s.favoritesCount(),
	}
	response.Tags = make([]string, 0)
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.c, tag}
		response.Tags = append(response.Tags, serializer.Response())
	}
	return response
}
