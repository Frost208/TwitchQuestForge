package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "twitch-rpg"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Character routes
		characters := v1.Group("/characters")
		{
			characterHandler := NewCharacterHandler()
			characters.POST("/", characterHandler.CreateCharacter)
			characters.GET("/:id", characterHandler.GetCharacter)
			characters.GET("/username/:username", characterHandler.GetCharacterByUsername)
			characters.PUT("/:id/stats", characterHandler.UpgradeStats)
			characters.PUT("/:id/equip", characterHandler.EquipItem)
			characters.DELETE("/:id/unequip/:slot", characterHandler.UnequipItem)
			characters.GET("/:id/inventory", characterHandler.GetInventory)
			characters.GET("/", characterHandler.GetAllCharacters)
		}

		// Item routes
		items := v1.Group("/items")
		{
			itemHandler := NewItemHandler()
			items.GET("/:id", itemHandler.GetItem)
			items.GET("/type/:type", itemHandler.GetItemsByType)
			items.GET("/random", itemHandler.GetRandomItems)
		}

		// Combat routes
		combat := v1.Group("/combat")
		{
			combatHandler := NewCombatHandler()
			combat.POST("/challenge", combatHandler.StartCombat)
			combat.GET("/history", combatHandler.GetCombatHistory)
		}

		// Game events routes (for OBS integration)
		events := v1.Group("/events")
		{
			eventHandler := NewEventHandler()
			events.GET("/latest", eventHandler.GetLatestEvents)
			events.PUT("/:id/trigger", eventHandler.MarkEventTriggered)
		}

		// Merchant routes
		merchant := v1.Group("/merchant")
		{
			merchantHandler := NewMerchantHandler()
			merchant.GET("/current", merchantHandler.GetCurrentEvent)
			merchant.POST("/create", merchantHandler.CreateMerchantEvent)
			merchant.POST("/purchase", merchantHandler.PurchaseItem)
		}
	}
}