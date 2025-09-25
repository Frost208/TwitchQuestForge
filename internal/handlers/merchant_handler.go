package handlers

import (
        "github.com/gin-gonic/gin"
        "net/http"
)

// MerchantHandler handles merchant-related HTTP requests
type MerchantHandler struct{}

// NewMerchantHandler creates a new merchant handler
func NewMerchantHandler() *MerchantHandler {
        return &MerchantHandler{}
}

// GetCurrentEvent retrieves current merchant event
func (mh *MerchantHandler) GetCurrentEvent(c *gin.Context) {
        c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

// CreateMerchantEvent creates a merchant event
func (mh *MerchantHandler) CreateMerchantEvent(c *gin.Context) {
        c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

// PurchaseItem handles item purchases
func (mh *MerchantHandler) PurchaseItem(c *gin.Context) {
        c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}