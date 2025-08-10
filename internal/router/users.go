package router

import (
	"OnlineShop/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"Test@gmail.com"`
	Password string `json:"password" binding:"required,min=8"`
}
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		Role:   role,
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
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, HTTPError{Message: "Access denied: role not found in token"})
			c.Abort()
			return
		}

		if userRole.(string) != "admin" {
			c.JSON(http.StatusForbidden, HTTPError{Message: "Access denied: requires admin privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// @Summary      Регистрация нового пользователя
// @Description  Создает новый аккаунт пользователя с email и паролем.
// @Tags         Пользователи (Auth)
// @Accept       json
// @Produce      json
// @Param        input  body      router.LoginInput  true  "Данные для регистрации"
// @Success      201    {object}  database.Customer     "Возвращает созданного пользователя"
// @Failure      400    {object}  router.HTTPError      "Ошибка валидации входных данных"
// @Failure      409    {object}  router.HTTPError      "Пользователь с таким email уже существует"
// @Failure      500    {object}  router.HTTPError      "Внутренняя ошибка сервера"
// @Router       /users/register [post]
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

// @Summary      Вход пользователя в систему
// @Description  Проверяет учетные данные и в случае успеха возвращает JWT токен.
// @Tags         Пользователи (Auth)
// @Accept       json
// @Produce      json
// @Param        input  body      router.LoginInput     true  "Учетные данные для входа"
// @Success      200    {object}  object{token=string}  "JWT токен"
// @Failure      400    {object}  router.HTTPError
// @Failure      401    {object}  router.HTTPError
// @Router       /users/login [post]
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

	token, err := GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary      Получить информацию о текущем пользователе
// @Description  Возвращает данные пользователя, аутентифицированного с помощью JWT токена.
// @Tags         Пользователи (Auth)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  database.Customer  "Данные текущего пользователя"
// @Failure      401  {object}  router.HTTPError   "Ошибка аутентификации"
// @Failure      404  {object}  router.HTTPError   "Пользователь из токена не найден в БД"
// @Router       /users/me [get]
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

// @Summary      Повысить пользователя до администратора
// @Description  Позволяет администратору назначить другого пользователя администратором.
// @Tags         Администрирование (Admin)
// @Produce      json
// @Param        id   path      int  true  "ID пользователя, которого нужно повысить"
// @Security     BearerAuth
// @Success      200  {object}  router.SuccessMessage "Сообщение об успешном повышении"
// @Failure      400  {object}  HTTPError           "Некорректный ID пользователя"
// @Failure      403  {object}  HTTPError           "Попытка повысить самого себя"
// @Failure      404  {object}  HTTPError           "Пользователь не найден"
// @Failure      409  {object}  HTTPError           "Пользователь уже является администратором"
// @Failure      500  {object}  HTTPError           "Внутренняя ошибка сервера"
// @Router       /users/{id}/promote [post]
func promoteUserToAdmin(c *gin.Context) {
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "Invalid user ID"})
		return
	}

	var user database.Customer
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, HTTPError{Message: "User not found"})
		return
	}

	if user.Role == "admin" {
		c.JSON(http.StatusConflict, HTTPError{Message: "User is already an admin"})
		return
	}

	if err := database.DB.Model(&user).Update("role", "admin").Error; err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, SuccessMessage{Message: "User successfully promoted to admin"})
}
