package user_grpc_service

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/gender"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	agent3 "github.com/Nistagram-Organization/nistagram-users/src/services/agent"
	registered_user2 "github.com/Nistagram-Organization/nistagram-users/src/services/registered_user"
	user3 "github.com/Nistagram-Organization/nistagram-users/src/services/user"
)

type userGrpcService struct {
	proto.UserServiceServer
	registeredUserService registered_user2.RegisteredUserService
	agentService          agent3.AgentService
	userService           user3.UserService
}

func NewUserGrpcService(agentService agent3.AgentService, registeredUserService registered_user2.RegisteredUserService, userService user3.UserService) proto.UserServiceServer {
	return &userGrpcService{
		proto.UnimplementedUserServiceServer{},
		registeredUserService,
		agentService,
		userService,
	}
}

func (s *userGrpcService) CreateUser(ctx context.Context, registrationRequest *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	userMessage := registrationRequest.GetRegistration()

	var id uint64
	if userMessage.Role == proto.Role_USER {
		var userEntity *user.User
		var err error

		userEntity = getUser(userMessage)
		registeredUser := &registered_user.RegisteredUser{
			User: *userEntity,
		}

		registeredUser, err = s.registeredUserService.Create(registeredUser)
		if err != nil {
			return nil, err
		}

		id = uint64(registeredUser.ID)
	} else {
		var agentEntity *agent.Agent
		var err error

		agentEntity = toAgent(userMessage)
		agentEntity, err = s.agentService.Create(agentEntity)

		if err != nil {
			return nil, err
		}

		id = uint64(agentEntity.ID)
	}

	res := proto.RegistrationResponse{Id: id}

	return &res, nil
}

func (s *userGrpcService) DeleteUser(ctx context.Context, deleteUserRequest *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if deleteUserRequest.Role == proto.Role_USER {
		err := s.registeredUserService.Delete(uint(deleteUserRequest.Id))
		if err != nil {
			return nil, err
		}
	} else {
		err := s.agentService.Delete(uint(deleteUserRequest.Id))
		if err != nil {
			return nil, err
		}
	}
	return &proto.DeleteUserResponse{Success: true}, nil
}

func toAgent(agentEntity *proto.UserMessage) *agent.Agent {
	userEntity := getUser(agentEntity)

	return &agent.Agent{
		User: *userEntity,
	}
}

func getUser(userEntity *proto.UserMessage) *user.User {
	return &user.User{
		Name:      userEntity.Username,
		Username:  userEntity.Username,
		FirstName: userEntity.Name,
		LastName:  userEntity.Surname,
		Phone:     userEntity.Phone,
		BirthDate: userEntity.BirthDate,
		Gender:    toGender(userEntity.Gender),
		Public:    userEntity.Public,
		Taggable:  userEntity.Taggable,
		Active:    true,
		Email:     userEntity.Email,
	}
}

func toGender(messageGender proto.UserMessage_Gender) gender.Gender {
	if messageGender == proto.UserMessage_MALE {
		return gender.Male
	} else {
		return gender.Female
	}
}
