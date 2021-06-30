package post_user_repository

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type PostUserRepository interface {
	Save(*post.PostUser) (*post.PostUser, rest_error.RestErr)
	Delete(uint) rest_error.RestErr
}

type postUserRepository struct {
	db *gorm.DB
}

func NewPostUserRepository(databaseClient datasources.DatabaseClient) PostUserRepository {
	return &postUserRepository{
		databaseClient.GetClient(),
	}
}

func (p *postUserRepository) Save(postUser *post.PostUser) (*post.PostUser, rest_error.RestErr) {
	if err := p.db.Create(&postUser).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to save post_user", err)
	}
	return postUser, nil
}

func (p *postUserRepository) Delete(id uint) rest_error.RestErr {
	postUser := post.PostUser{
		ID: id,
	}
	if err := p.db.Where("id = ?", id).Delete(postUser).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete post_user", err)
	}
	return nil
}
