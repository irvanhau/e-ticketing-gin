package server

import (
	"e-ticketing-gin/configs"
	"e-ticketing-gin/utils/database"
	"e-ticketing-gin/utils/database/seeds"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	g *gin.Engine
	c *configs.ProgramConfig
}

func (s *Server) RunServer() {
	s.g.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)
	}))

	s.g.Use(gin.Recovery())
	err := s.g.Run(":8080")
	if err != nil {
		logrus.Error("Failed to start server : ", err.Error())
		return
	}
}

func (s *Server) MigrateDB() {
	db := database.InitDB(s.c)
	database.Migrate(db)
}

func (s *Server) SeederDB() {
	db := database.InitDB(s.c)
	for _, seed := range seeds.All() {
		if err := seed.Run(db); err != nil {
			fmt.Printf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}

func InitServer(g *gin.Engine, c *configs.ProgramConfig) *Server {
	return &Server{
		g: g,
		c: c,
	}
}
