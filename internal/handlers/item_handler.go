package handlers

import (
        "net/http"
        "strconv"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/services"

        "github.com/gin-gonic/gin"
)

// ItemHandler handles item-related HTTP requests
type ItemHandler struct {
        itemService *services.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler() *ItemHandler {
        return &ItemHandler{
                itemService: services.NewItemService(),
        }
}

// GetItem retrieves an item by ID
func (ih *ItemHandler) GetItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
                return
        }

        item, err := ih.itemService.GetItemByID(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        if item == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
                return
        }

        c.JSON(http.StatusOK, item)
}

// GetItemsByType retrieves items by type
func (ih *ItemHandler) GetItemsByType(c *gin.Context) {
        itemType := c.Param("type")

        // Validate item type
        validTypes := []models.ItemType{models.ItemTypeBoots, models.ItemTypePants, models.ItemTypeArmor, models.ItemTypeHelmet, models.ItemTypeRing, models.ItemTypeChain}
        validType := false
        for _, valid := range validTypes {
                if models.ItemType(itemType) == valid {
                        validType = true
                        break
                }
        }

        if !validType {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
                return
        }

        // Parse optional pagination parameters
        limit := 50 // default
        offset := 0 // default

        if limitStr := c.Query("limit"); limitStr != "" {
                if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
                        limit = l
                }
        }

        if offsetStr := c.Query("offset"); offsetStr != "" {
                if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
                        offset = o
                }
        }

        items, err := ih.itemService.GetItemsByType(models.ItemType(itemType), limit, offset)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"items": items, "count": len(items)})
}

// GetRandomItems retrieves random items
func (ih *ItemHandler) GetRandomItems(c *gin.Context) {
        // Parse optional parameters
        count := 5 // default
        isSpecial := false // default

        if countStr := c.Query("count"); countStr != "" {
                if c, err := strconv.Atoi(countStr); err == nil && c > 0 && c <= 20 {
                        count = c
                }
        }

        if specialStr := c.Query("special"); specialStr == "true" {
                isSpecial = true
        }

        items, err := ih.itemService.GetRandomItems(count, isSpecial)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"items": items, "count": len(items)})
}