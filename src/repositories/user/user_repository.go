package user

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type UserRepository interface {
	Delete(uint) rest_error.RestErr
	GetByEmail(string) (*user.User, rest_error.RestErr)
	Update(*user.User) (*user.User, rest_error.RestErr)
	DeleteFavorite(uint, uint) rest_error.RestErr
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(databaseClient datasources.DatabaseClient) UserRepository {
	return &userRepository{
		databaseClient.GetClient(),
	}
}

func (r *userRepository) Create(user *registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create user", err)
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) rest_error.RestErr {
	if err := r.db.Where("owner_id = ?", id).Delete(&user.User{}).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete user", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(email string) (*user.User, rest_error.RestErr) {
	userEntity := user.User{
		Email: email,
	}
	if err := r.db.Where("email = ?", email).Preload("Favorites").First(&userEntity).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))
	}

	return &userEntity, nil
}

func (r *userRepository) Update(user *user.User) (*user.User, rest_error.RestErr) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to update user", err)
	}
	return user, nil
}

func (r *userRepository) DeleteFavorite(userId uint, postUserId uint) rest_error.RestErr {
	tx := r.db.Exec(fmt.Sprintf("delete from favorites where user_id=%d & post_user_id=%d", userId, postUserId))

	if tx.Error != nil {
		return rest_error.NewInternalServerError("Error when trying to delete post from favorites", tx.Error)
	}

	return nil
}
