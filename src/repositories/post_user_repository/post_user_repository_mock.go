package post_user_repository

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type PostUserRepositoryMock struct {
	mock.Mock
}

func (p *PostUserRepositoryMock) Save(user *post.PostUser) (*post.PostUser, rest_error.RestErr) {
	panic("implement me")
}

func (p *PostUserRepositoryMock) Delete(u uint) rest_error.RestErr {
	panic("implement me")
}

