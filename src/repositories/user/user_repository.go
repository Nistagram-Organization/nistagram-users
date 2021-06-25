package user

import (
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*user.User) (*user.User, rest_error.RestErr)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(databaseClient datasources.DatabaseClient) UserRepository {
	return &userRepository{
		databaseClient.GetClient(),
	}
}

func (r *userRepository) Create(user *user.User) (*user.User, rest_error.RestErr) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create user", err)
	}
	return user, nil
}
