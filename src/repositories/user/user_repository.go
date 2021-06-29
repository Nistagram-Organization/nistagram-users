package user

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*user.User, rest_error.RestErr)
	Create(*user.User) (*user.User, rest_error.RestErr)
	Edit(*user.User) (*user.User, rest_error.RestErr)
	Delete(uint) rest_error.RestErr
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(databaseClient datasources.DatabaseClient) UserRepository {
	return &userRepository{
		databaseClient.GetClient(),
	}
}

func (r *userRepository) GetByEmail(email string) (*user.User, rest_error.RestErr) {
	user := user.User{
		Email: email,
	}
	if err := r.db.Take(&user, user.Email).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get user with email %s", user.Email))
	}
	return &user, nil
}

func (r *userRepository) Create(user *user.User) (*user.User, rest_error.RestErr) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create user", err)
	}
	return user, nil
}

func (r *userRepository) Edit(user *user.User) (*user.User, rest_error.RestErr) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to edit user", err)
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) rest_error.RestErr {
	if err := r.db.Delete(&user.User{}, id).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete agent", err)
	}
	return nil
}
