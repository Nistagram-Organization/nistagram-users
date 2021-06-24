package application

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/campaign"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-users/src/controllers/ping"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	"github.com/gin-gonic/gin"
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

	pingController := ping.NewPingController()

	router.GET("/ping", pingController.Ping)

	router.Run(":8084")
}
