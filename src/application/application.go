package application

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	usercontroller "github.com/Nistagram-Organization/nistagram-users/src/controllers/user"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	agent2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/post_user_repository"
	registered_user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	agent3 "github.com/Nistagram-Organization/nistagram-users/src/services/agent"
	registered_user3 "github.com/Nistagram-Organization/nistagram-users/src/services/registered_user"
	user3 "github.com/Nistagram-Organization/nistagram-users/src/services/user"
	"github.com/Nistagram-Organization/nistagram-users/src/services/user_grpc_service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	router = gin.Default()
)

func StartApplication() {
	router.Use(cors.Default())

	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := database.Migrate(
		&user.User{},
		&registered_user.RegisteredUser{},
		&agent.Agent{},
		&post.PostUser{},
	); err != nil {
		panic(err)
	}

	port := ":8084"
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1%s", port))
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	grpcS := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcS, user_grpc_service.NewUserGrpcService(
		agent3.NewAgentService(
			agent2.NewAgentRepository(database),
			user2.NewUserRepository(database),
		),
		registered_user3.NewRegisteredUserService(
			registered_user2.NewRegisteredUserRepository(database),
			user2.NewUserRepository(database),
		),
		user3.NewUserService(
			user2.NewUserRepository(database),
			registered_user2.NewRegisteredUserRepository(database),
			post_user_repository.NewPostUserRepository(database),
		),
	))

	userController := usercontroller.NewUserController(
		user3.NewUserService(
			user2.NewUserRepository(database),
			registered_user2.NewRegisteredUserRepository(database),
			post_user_repository.NewPostUserRepository(database),
		),
	)
	router.POST("/users/favorites", userController.AddPostToFavorites)
	router.DELETE("/users/favorites", userController.RemovePostFromFavorites)

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
