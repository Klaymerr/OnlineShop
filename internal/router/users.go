package router

import (
	"OnlineShop/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, HTTPError{Message: "missing or invalid token"})
			c.Abort()
			return
		}

		tokenStr = tokenStr[7:]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, HTTPError{Message: "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func registerUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
		return
	}

	var existingUser database.Customer
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, HTTPError{Message: "user with this email already exists"}) // 409 Conflict
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "failed to hash password"})
		return
	}

	newUser := database.Customer{
		Email:            input.Email,
		PasswordHash:     string(hashedPassword),
		RegistrationDate: time.Now(),
	}

	if result := database.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func loginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "invalid request body"})
		return
	}

	var user database.Customer
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, HTTPError{Message: "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, HTTPError{Message: "invalid email or password"})
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SayHello(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, HTTPError{Message: "user ID not found in context"})
		return
	}

	id := userID.(uint)

	var user database.Customer
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, HTTPError{Message: "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
