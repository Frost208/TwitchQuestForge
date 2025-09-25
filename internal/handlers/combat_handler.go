package handlers

import (
        "net/http"
        "strconv"
        "twitch-rpg/internal/services"

        "github.com/gin-gonic/gin"
)

// CombatHandler handles combat-related HTTP requests
type CombatHandler struct {
        combatService *services.CombatService
}

// NewCombatHandler creates a new combat handler
func NewCombatHandler() *CombatHandler {
        return &CombatHandler{
                combatService: services.NewCombatService(),
        }
}

// StartCombat starts a combat encounter
func (ch *CombatHandler) StartCombat(c *gin.Context) {
        var req struct {
                AttackerID int `json:"attacker_id" binding:"required"`
                DefenderID int `json:"defender_id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        result, err := ch.combatService.StartCombat(req.AttackerID, req.DefenderID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, result)
}

// GetCombatHistory retrieves combat history
func (ch *CombatHandler) GetCombatHistory(c *gin.Context) {
        limit := 20 // default
        if limitStr := c.Query("limit"); limitStr != "" {
                if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
                        limit = l
                }
        }

        combatHistory, err := ch.combatService.GetCombatHistory(limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"combat_history": combatHistory, "count": len(combatHistory)})
}