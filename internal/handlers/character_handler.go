package handlers

import (
        "net/http"
        "strconv"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/services"

        "github.com/gin-gonic/gin"
)

// CharacterHandler handles character-related HTTP requests
type CharacterHandler struct {
        characterService *services.CharacterService
}

// NewCharacterHandler creates a new character handler
func NewCharacterHandler() *CharacterHandler {
        return &CharacterHandler{
                characterService: services.NewCharacterService(),
        }
}

// CreateCharacter creates a new character
func (ch *CharacterHandler) CreateCharacter(c *gin.Context) {
        var req models.CharacterCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        character, err := ch.characterService.CreateCharacter(req.Username, req.TwitchUserID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusCreated, character)
}

// GetCharacter retrieves a character by ID
func (ch *CharacterHandler) GetCharacter(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
                return
        }

        character, err := ch.characterService.GetCharacterByID(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        if character == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
                return
        }

        c.JSON(http.StatusOK, character)
}

// GetCharacterByUsername retrieves a character by username
func (ch *CharacterHandler) GetCharacterByUsername(c *gin.Context) {
        username := c.Param("username")

        character, err := ch.characterService.GetCharacterByUsername(username)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        if character == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
                return
        }

        c.JSON(http.StatusOK, character)
}

// UpgradeStats upgrades character stats
func (ch *CharacterHandler) UpgradeStats(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
                return
        }

        var req models.StatUpgradeRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        character, err := ch.characterService.UpgradeCharacterStat(id, req.StatType, req.ChannelPoints)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, character)
}

// EquipItem equips an item to character
func (ch *CharacterHandler) EquipItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
                return
        }

        var req models.EquipItemRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        err = ch.characterService.EquipItem(id, req.ItemID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Item equipped successfully"})
}

// UnequipItem unequips an item from character
func (ch *CharacterHandler) UnequipItem(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
                return
        }

        slotStr := c.Param("slot")
        slotType := models.ItemType(slotStr)

        // Validate slot type
        validSlots := []models.ItemType{models.ItemTypeBoots, models.ItemTypePants, models.ItemTypeArmor, models.ItemTypeHelmet, models.ItemTypeRing, models.ItemTypeChain}
        validSlot := false
        for _, valid := range validSlots {
                if slotType == valid {
                        validSlot = true
                        break
                }
        }

        if !validSlot {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot type"})
                return
        }

        err = ch.characterService.UnequipItem(id, slotType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Item unequipped successfully"})
}

// GetInventory gets character inventory
func (ch *CharacterHandler) GetInventory(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
                return
        }

        inventory, err := ch.characterService.GetCharacterInventory(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"inventory": inventory})
}

// GetAllCharacters gets all characters
func (ch *CharacterHandler) GetAllCharacters(c *gin.Context) {
        characters, err := ch.characterService.GetAllCharacters()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, characters)
}