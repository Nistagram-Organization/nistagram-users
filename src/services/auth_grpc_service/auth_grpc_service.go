package auth_grpc_service

import (
	"context"
	agent2 "github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/gender"
	user2 "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
)

type authGrpcService struct {
	proto.AuthServiceServer
	agentRepository agent.AgentRepository
	userRepository  user.UserRepository
}

func NewAuthGrpcService(agentRepository agent.AgentRepository, userRepository user.UserRepository) proto.AuthServiceServer {
	return &authGrpcService{
		proto.UnimplementedAuthServiceServer{},
		agentRepository,
		userRepository,
	}
}

func (s *authGrpcService) Register(ctx context.Context, registrationRequest *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	userMessage := registrationRequest.GetRegistration()

	if userMessage.Role == proto.Role_USER {
		user := toUser(userMessage)

		if _, err := s.userRepository.Create(&user); err != nil {
			return nil, err
		}
	} else {
		agent := toAgent(userMessage)

		if _, err := s.agentRepository.Create(&agent); err != nil {
			return nil, err
		}
	}
	res := proto.RegistrationResponse{Success: true}

	return &res, nil
}

func toAgent(agent *proto.UserMessage) agent2.Agent {
	user := toUser(agent)
	user.Active = false

	return agent2.Agent{
		User: user,
	}

}

func toUser(user *proto.UserMessage) user2.User {
	return user2.User{
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
