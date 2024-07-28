package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateReview godoc
// @Summary Create a new review
// @Description Create a new review
// @Tags reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param review body models.Review true "Review"
// @Success 200 {object} models.Review
// @Router /reviews [post]
func CreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// GetReviews godoc
// @Summary Get reviews by phone ID
// @Description Get reviews by phone ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param phone_id path int true "Phone ID"
// @Success 200 {object} []models.Review
// @Router /reviews/{phone_id} [get]
func GetReviews(c *gin.Context) {
	phoneID := c.Param("phone_id")
	var reviews []models.Review
	if err := config.DB.Where("phone_id = ?", phoneID).Preload("Comments").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// GetAllReviews godoc
// @Summary Get all reviews
// @Description Get all reviews without authentication
// @Tags reviews
// @Accept json
// @Produce json
// @Success 200 {object} []models.Review
// @Router /reviews [get]
func GetAllReviews(c *gin.Context) {
	var reviews []models.Review
	if err := config.DB.Preload("Comments").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// UpdateReview godoc
// @Summary Update a review
// @Description Update a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "Review ID"
// @Param review body models.Review true "Review"
// @Success 200 {object} models.Review
// @Router /reviews/{id} [put]
func UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")

	var updatedReview models.Review
	if err := c.ShouldBindJSON(&updatedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingReview models.Review
	result := config.DB.First(&existingReview, reviewID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	existingReview.Rating = updatedReview.Rating
	existingReview.Content = updatedReview.Content

	if err := config.DB.Save(&existingReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingReview)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "Review ID"
// @Success 200 {object} map[string]string
// @Router /reviews/{id} [delete]
func DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")

	var existingReview models.Review
	result := config.DB.First(&existingReview, reviewID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if err := config.DB.Delete(&existingReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}

// GetReviewByID godoc
// @Summary Get a review by ID
// @Description Get a review by ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} models.Review
// @Router /reviews/{id} [get]
func GetReviewByID(c *gin.Context) {
	reviewID := c.Param("id")
	var review models.Review

	// Retrieve the review from the database
	if err := config.DB.Preload("Comments").First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the review details
	c.JSON(http.StatusOK, review)
}
