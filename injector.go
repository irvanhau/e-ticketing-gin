//go:build wireinject
// +build wireinject

package main

import (
	"e-ticketing-gin/configs"
	"e-ticketing-gin/features/users"
	userData "e-ticketing-gin/features/users/data"
	userHandler "e-ticketing-gin/features/users/handler"
	userService "e-ticketing-gin/features/users/service"
	"e-ticketing-gin/helper/email"
	"e-ticketing-gin/helper/enkrip"
	"e-ticketing-gin/helper/jwt"
	"e-ticketing-gin/routes"
	"e-ticketing-gin/server"
	"e-ticketing-gin/utils/database"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	userData.New,
	wire.Bind(new(users.UserDataInterface), new(*userData.UserData)),

	userService.New,
	wire.Bind(new(users.UserServiceInterface), new(*userService.UserService)),

	userHandler.NewHandler,
	wire.Bind(new(users.UserHandlerInterface), new(*userHandler.UserHandler)),
)

func InitializedServer() *server.Server {
	wire.Build(
		configs.InitConfig,
		database.InitDB,
		enkrip.New,
		email.NewEmail,
		jwt.NewJWT,
		//JANGAN DIUBAH

		userSet,

		// JANGAN DIUBAH
		routes.NewRoute,
		server.InitServer,
	)
	return nil
}
