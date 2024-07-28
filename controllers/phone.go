package controllers

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPhones godoc
// @Summary Get all phones
// @Description Get all phones
// @Tags phones
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Phone
// @Router /phones [get]
func GetPhones(c *gin.Context) {
	var phones []models.Phone

	if err := config.DB.Preload("Features").Preload("Reviews").Find(&phones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, phones)
}

// CreatePhone godoc
// @Summary Create a new phone
// @Description Create a new phone
// @Tags phones
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone body models.PhoneRequest true "Phone"
// @Success 200 {object} map[string]string
// @Router /phones [post]
func CreatePhone(c *gin.Context) {
	var input models.Phone
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "phone created successfully"})
}

// GetPhoneByID godoc
// @Summary Get a phone by ID
// @Description Get a phone by ID
// @Tags phones
// @Accept  json
// @Produce  json
// @Param phone_id path int true "Phone ID"
// @Success 200 {object} models.Phone
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /phones/{phone_id} [get]
func GetPhoneByID(c *gin.Context) {
	phoneID := c.Param("phone_id")

	var phone models.Phone
	if err := config.DB.Preload("Features").Preload("Reviews").First(&phone, phoneID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "phone not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, phone)
}

// UpdatePhone godoc
// @Summary Update a phone
// @Description Update a phone
// @Tags phones
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone_id path int true "Phone ID"
// @Param phone body models.PhoneRequest true "Phone"
// @Success 200 {object} models.Phone
// @Router /phones/{phone_id} [put]
func UpdatePhone(c *gin.Context) {
	phoneID := c.Param("phone_id")

	var input models.Phone
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert phone_id to uint
	id, err := strconv.ParseUint(phoneID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone ID"})
		return
	}

	input.ID = uint(id)

	if err := config.DB.Save(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DeletePhone godoc
// @Summary Delete a phone by ID
// @Description Delete a phone by ID
// @Tags phones
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Authorization header"
// @Security ApiKeyAuth
// @Param phone_id path int true "Phone ID"
// @Success 200 {string} string "Phone deleted successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /phones/{phone_id} [delete]
func DeletePhone(c *gin.Context) {
	phoneID := c.Param("phone_id")

	// Convert phone_id to uint
	id, err := strconv.ParseUint(phoneID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone ID"})
		return
	}

	var phone models.Phone
	result := config.DB.First(&phone, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "phone not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if err := config.DB.Delete(&phone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "phone deleted successfully"})
}
