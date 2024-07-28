package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"backend-vercel-phone-review/utils"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateFeature godoc
// @Summary Create a new feature
// @Description Create a new feature
// @Tags features
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone_id path int true "Phone ID"
// @Param feature body models.Feature true "Feature"
// @Success 200 {object} models.Feature
// @Router /phones/{phone_id}/features [post]
func CreateFeature(c *gin.Context) {
	phoneID := c.Param("phone_id")
	var feature models.Feature

	if err := c.ShouldBindJSON(&feature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set PhoneID for the feature
	feature.PhoneID = utils.StringToUint(phoneID)

	if err := config.DB.Create(&feature).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feature)
}

// UpdateFeature godoc
// @Summary Update a feature of a phone
// @Description Update a feature of a phone
// @Tags features
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone_id path int true "Phone ID"
// @Param feature_id path int true "Feature ID"
// @Param feature body models.Feature true "Feature"
// @Success 200 {object} models.Feature
// @Router /phones/{phone_id}/features/{feature_id} [put]
func UpdateFeature(c *gin.Context) {
	phoneID, err := strconv.Atoi(c.Param("phone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone ID"})
		return
	}

	featureID, err := strconv.Atoi(c.Param("feature_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feature ID"})
		return
	}

	var feature models.Feature
	if err := c.ShouldBindJSON(&feature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingFeature models.Feature
	result := config.DB.First(&existingFeature, featureID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feature not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	// Check if the feature belongs to the correct phone
	if existingFeature.PhoneID != uint(phoneID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feature does not belong to the specified phone"})
		return
	}

	existingFeature.Name = feature.Name
	existingFeature.Details = feature.Details

	if err := config.DB.Save(&existingFeature).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feature updated successfully"})
}

// DeleteFeature godoc
// @Summary Delete a feature of a phone
// @Description Delete a feature of a phone
// @Tags features
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone_id path int true "Phone ID"
// @Param feature_id path int true "Feature ID"
// @Success 200 {object} models.Feature
// @Router /phones/{phone_id}/features/{feature_id} [delete]
func DeleteFeature(c *gin.Context) {
	phoneID, err := strconv.Atoi(c.Param("phone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone ID"})
		return
	}

	featureID, err := strconv.Atoi(c.Param("feature_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feature ID"})
		return
	}

	var existingFeature models.Feature
	result := config.DB.First(&existingFeature, featureID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feature not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	// Check if the feature belongs to the correct phone
	if existingFeature.PhoneID != uint(phoneID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feature does not belong to the specified phone"})
		return
	}

	if err := config.DB.Delete(&existingFeature).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feature deleted successfully"})
}
