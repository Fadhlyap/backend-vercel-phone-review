package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param comment body models.Comment true "Comment"
// @Success 200 {object} models.Comment
// @Router /comments [post]
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetComments godoc
// @Summary Get comments by review ID
// @Description Get comments by review ID
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param review_id path int true "Review ID"
// @Success 200 {object} []models.Comment
// @Router /comments/{review_id} [get]
func GetComments(c *gin.Context) {
	reviewID := c.Param("review_id")

	var comments []models.Comment
	if err := config.DB.Where("review_id = ?", reviewID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "Comment ID"
// @Param comment body models.Comment true "Comment"
// @Success 200 {object} map[string]string
// @Router /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	commentID := c.Param("id")

	var updatedComment models.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingComment models.Comment
	result := config.DB.First(&existingComment, commentID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	existingComment.Content = updatedComment.Content

	if err := config.DB.Save(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingComment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	var existingComment models.Comment
	result := config.DB.First(&existingComment, commentID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if err := config.DB.Delete(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
