package user_grpc_service

import (
	"context"
	agent2 "github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/gender"
	user2 "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	agent3 "github.com/Nistagram-Organization/nistagram-users/src/services/agent"
	user3 "github.com/Nistagram-Organization/nistagram-users/src/services/user"
)

type userGrpcService struct {
	proto.UserServiceServer
	agentService agent3.AgentService
	userService  user3.UserService
}

func NewUserGrpcService(agentService agent3.AgentService, userService user3.UserService) proto.UserServiceServer {
	return &userGrpcService{
		proto.UnimplementedUserServiceServer{},
		agentService,
		userService,
	}
}

func (s *userGrpcService) CreateUser(ctx context.Context, registrationRequest *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	userMessage := registrationRequest.GetRegistration()

	var id uint64
	if userMessage.Role == proto.Role_USER {
		var user *user2.User
		var err error

		user = toUser(userMessage)
		user, err = s.userService.Create(user)
		if err != nil {
			return nil, err
		}

		id = uint64(user.ID)
	} else {
		var agent *agent2.Agent
		var err error

		agent = toAgent(userMessage)
		agent, err = s.agentService.Create(agent)

		if err != nil {
			return nil, err
		}

		id = uint64(agent.ID)
	}

	res := proto.RegistrationResponse{Id: id}

	return &res, nil
}

func (s *userGrpcService) DeleteUser(ctx context.Context, deleteUserRequest *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if deleteUserRequest.Role == proto.Role_USER {
		err := s.userService.Delete(uint(deleteUserRequest.Id))
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

func toAgent(agent *proto.UserMessage) *agent2.Agent {
	user := toUser(agent)
	user.Active = false

	return &agent2.Agent{
		User: *user,
	}

}

func toUser(user *proto.UserMessage) *user2.User {
	return &user2.User{
		Username:  user.Username,
		Name:      user.Name,
		Surname:   user.Surname,
		Phone:     user.Password,
		BirthDate: user.BirthDate,
		Gender:    toGender(user.Gender),
		Public:    user.Public,
		Taggable:  user.Taggable,
		Active:    true,
		Email:     user.Email,
	}
}

func toGender(messageGender proto.UserMessage_Gender) gender.Gender {
	if messageGender == proto.UserMessage_MALE {
		return gender.Male
	} else {
		return gender.Female
	}
}
