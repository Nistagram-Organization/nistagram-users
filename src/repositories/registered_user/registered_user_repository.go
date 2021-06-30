package registered_user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type RegisteredUserRepository interface {
	Create(*registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr)
	Delete(uint) rest_error.RestErr
}

type registeredUserRepository struct {
	db *gorm.DB
}

func NewRegisteredUserRepository(databaseClient datasources.DatabaseClient) RegisteredUserRepository {
	return &registeredUserRepository{
		databaseClient.GetClient(),
	}
}

func (r *registeredUserRepository) Create(user *registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create registered user", err)
	}
	return user, nil
}

func (r *registeredUserRepository) Delete(id uint) rest_error.RestErr {
	if err := r.db.Delete(&registered_user.RegisteredUser{}, id).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete registered user", err)
	}
	return nil
}
