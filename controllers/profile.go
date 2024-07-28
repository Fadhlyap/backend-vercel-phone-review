package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Profile  struct {
		FullName string `json:"full_name"`
		Bio      string `json:"bio"`
	} `json:"profile,omitempty"`
	Reviews []ReviewResponse `json:"reviews,omitempty"`
}

type ReviewResponse struct {
	ID      uint   `json:"id"`
	PhoneID uint   `json:"phone_id"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")

	if err := config.DB.Preload("Profile").Preload("Reviews").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	if user.Profile.ID != 0 {
		userResponse.Profile.FullName = user.Profile.FullName
		userResponse.Profile.Bio = user.Profile.Bio
	}

	for _, review := range user.Reviews {
		userResponse.Reviews = append(userResponse.Reviews, ReviewResponse{
			ID:      review.ID,
			PhoneID: review.PhoneID,
			Rating:  review.Rating,
			Content: review.Content,
		})
	}

	c.JSON(http.StatusOK, userResponse)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile
// @Tags users
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param profile body models.Profile true "Profile"
// @Success 200 {object} map[string]string
// @Router /users/{id}/profile [put]
func UpdateProfile(c *gin.Context) {
	var input struct {
		Bio      string `json:"bio"`
		FullName string `json:"full_name"`
	}
	userID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var profile models.Profile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newProfile := models.Profile{
				UserID:   user.ID,
				Bio:      input.Bio,
				FullName: input.FullName,
			}
			if err := config.DB.Create(&newProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if profile.UserID != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID in profile does not match requested user ID"})
			return
		}

		updateFields := map[string]interface{}{
			"bio":       input.Bio,
			"full_name": input.FullName,
		}
		if err := config.DB.Model(&profile).Updates(updateFields).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
