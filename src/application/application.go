package application

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/Nistagram-Organization/nistagram-users/src/controllers/ping"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	agent2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/Nistagram-Organization/nistagram-users/src/services/auth_grpc_service"
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
	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := database.Migrate(
		&user.User{},
		&agent.Agent{},
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
	proto.RegisterAuthServiceServer(grpcS, auth_grpc_service.NewAuthGrpcService(
		agent2.NewAgentRepository(database),
		user2.NewUserRepository(database),
	))

	pingController := ping.NewPingController()
	router.GET("/ping", pingController.Ping)

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
