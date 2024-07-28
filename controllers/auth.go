package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"backend-vercel-phone-review/utils"
	"net/http"
	"os"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegistRequest true "Regist"
// @Success 200 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Password = utils.HashPassword(input.Password)

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku selama 24 jam
	claims := &Claims{
		Username: input.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param password body models.ChangePasswordRequest true "Password"
// @Success 200 {object} map[string]string
// @Router /auth/change-password/{id} [put]
func ChangePassword(c *gin.Context) {
	userID := c.Param("id")
	var request models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if !utils.CheckPassword(request.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	newPassword := utils.HashPassword(request.NewPassword)
	if err := config.DB.Model(&user).Update("password", newPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
