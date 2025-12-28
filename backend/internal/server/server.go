package server

import (
	"auth-jwt/backend/internal/database"
	"auth-jwt/backend/internal/repositories"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")

	db := database.New()
	userRepository := repositories.NewUsersRepository(db)

	server := &Server{
		usersRepository: userRepository,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      server.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	authHandler := newAuthHandler(s.usersRepository)
	userHandler := newUserHandler(s.usersRepository)

	// Setup routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
		authRoutes.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
	}

	protected := r.Group("/api")
	{
		protected.GET("/user", userHandler.GetUser)
	}

	return r
}
