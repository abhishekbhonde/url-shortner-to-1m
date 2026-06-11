package handlers

import (
	"crypto/rand"
	"net/http"
	"time"
	"c/database"
	"c/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const codeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// fallback to time-based seed (extremely unlikely to hit)
		now := time.Now().UnixNano()
		for i := range b {
			b[i] = byte(now % int64(len(codeAlphabet)))
			now /= int64(len(codeAlphabet))
		}
	}
	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(b)
}

type shortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

// CreateShortURL creates and stores a shortened URL.
func CreateShortURL(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique short code (retry on collision)
	var code string
	for i := 0; i < 5; i++ {
		code = generateCode(6)
		var existing models.ShortURL
		if err := database.DB.Where("short_code = ?", code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		code = "" // collision; try again
	}
	if code == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate unique short code"})
		return
	}

	url := models.ShortURL{
		OriginalURL: req.OriginalURL,
		ShortCode:   code,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}

	short := fmt.Sprintf("http://%s/%s", c.Request.Host, url.ShortCode)
	c.JSON(http.StatusCreated, gin.H{"short_url": short, "code": url.ShortCode})
}

// RedirectShortURL redirects a short code to its original URL and increments click count.
func RedirectShortURL(c *gin.Context) {
	code := c.Param("code")
	var url models.ShortURL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// increment click counter efficiently
	database.DB.Model(&url).Where("id = ?", url.ID).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
	c.Redirect(http.StatusFound, url.OriginalURL)
}

// GetStats returns stats for a short code.
func GetStats(c *gin.Context) {
	code := c.Param("code")
	var url models.ShortURL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, url)
}

// DeleteShortURL removes a short URL by code.
func DeleteShortURL(c *gin.Context) {
	code := c.Param("code")
	var url models.ShortURL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	database.DB.Delete(&url)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
