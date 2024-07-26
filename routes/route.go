package routes

import (
	"e-ticketing-gin/features/users"
	"e-ticketing-gin/helper"
	"e-ticketing-gin/helper/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func NewRoute(uh users.UserHandlerInterface) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	jwtAuth := authMiddleware()

	api := router.Group("/api/v1")

	// Route Authentication
	api.POST("/register", uh.Register)
	api.POST("/login", uh.Login)
	api.POST("/forget-password", uh.ForgetPasswordWeb)
	api.POST("/reset-password", uh.ResetPassword)
	api.POST("/refresh-token", jwtAuth, uh.RefreshToken)
	api.POST("/verification", uh.UserVerification)

	// Route Profile
	api.GET("/profile", jwtAuth, uh.Profile)
	api.PUT("/profile/update", jwtAuth, uh.UpdateProfile)

	// Route User - Admin
	api.GET("/user", jwtAuth, uh.GetUsers)
	api.GET("/user/:id/activate", jwtAuth, uh.ActivateUser)
	api.GET("/user/:id/deactivate", jwtAuth, uh.DeactivateUser)
	api.GET("/user/dashboard", jwtAuth, uh.UserDashboard)

	return router
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.FormatResponse("Unauthorized", nil))
			return
		}
	}
}
