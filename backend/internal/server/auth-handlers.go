package server

import (
	"auth-jwt/backend/internal/auth"
	"auth-jwt/backend/internal/models"
	"auth-jwt/backend/internal/repositories"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepository repositories.UsersRepositoryInterface
}

func newAuthHandler(userRepository repositories.UsersRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		userRepository: userRepository,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the user already exists
	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// create a new user
	user, err = models.NewUser(req.Username, req.Email, hashedPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error creating user: %s", err.Error())})
		return
	}

	// save the user to the database
	err = h.userRepository.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// hash password & compare with the db password hash
	if !auth.ChackPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// save the token in cookies
	c.SetCookie(
		"access_token",
		accessToken,
		15*60,
		"/",
		"localhost",
		false,
		true, // não acessível via JS
	)
}
