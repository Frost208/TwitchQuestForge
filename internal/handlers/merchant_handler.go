package handlers

import (
        "net/http"
        "twitch-rpg/internal/services"

        "github.com/gin-gonic/gin"
)

// MerchantHandler handles merchant-related HTTP requests
type MerchantHandler struct {
        merchantService *services.MerchantService
}

// NewMerchantHandler creates a new merchant handler
func NewMerchantHandler() *MerchantHandler {
        return &MerchantHandler{
                merchantService: services.NewMerchantService(),
        }
}

// GetCurrentEvent retrieves current merchant event
func (mh *MerchantHandler) GetCurrentEvent(c *gin.Context) {
        event, err := mh.merchantService.GetCurrentEvent()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        if event == nil {
                c.JSON(http.StatusNotFound, gin.H{"message": "No active merchant event"})
                return
        }

        c.JSON(http.StatusOK, event)
}

// CreateMerchantEvent creates a merchant event
func (mh *MerchantHandler) CreateMerchantEvent(c *gin.Context) {
        var req struct {
                Title           string `json:"title" binding:"required"`
                Description     string `json:"description" binding:"required"`
                DurationMinutes int    `json:"duration_minutes" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        event, err := mh.merchantService.CreateMerchantEvent(req.Title, req.Description, req.DurationMinutes)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusCreated, event)
}

// PurchaseItem handles item purchases
func (mh *MerchantHandler) PurchaseItem(c *gin.Context) {
        var req struct {
                CharacterID         int `json:"character_id" binding:"required"`
                MerchantEventItemID int `json:"merchant_event_item_id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        err := mh.merchantService.PurchaseItem(req.CharacterID, req.MerchantEventItemID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Item purchased successfully"})
}