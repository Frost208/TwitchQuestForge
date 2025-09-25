package services

import (
        "fmt"
        "math/rand"
        "time"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/storage"
)

// MerchantService handles merchant-related operations
type MerchantService struct{}

// NewMerchantService creates a new merchant service
func NewMerchantService() *MerchantService {
        return &MerchantService{}
}

// GetCurrentEvent retrieves the current active merchant event
func (ms *MerchantService) GetCurrentEvent() (*models.MerchantEvent, error) {
        if database.DB == nil {
                return storage.Memory.GetCurrentMerchant()
        }

        query := `
                SELECT me.id, me.event_type, me.available_items, me.start_time, me.end_time, me.is_active
                FROM merchant_events me
                WHERE me.is_active = true AND (me.end_time IS NULL OR me.end_time > NOW())
                ORDER BY me.start_time DESC
                LIMIT 1`

        rows, err := database.DB.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to get current merchant event: %v", err)
        }
        defer rows.Close()

        var merchantEvent *models.MerchantEvent

        if rows.Next() {
                merchantEvent = &models.MerchantEvent{}
                err := rows.Scan(
                        &merchantEvent.ID, &merchantEvent.EventType, &merchantEvent.AvailableItems,
                        &merchantEvent.StartTime, &merchantEvent.EndTime, &merchantEvent.IsActive,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan merchant event: %v", err)
                }
        }

        return merchantEvent, nil
}

// CreateMerchantEvent creates a new merchant event with random items
func (ms *MerchantService) CreateMerchantEvent(title, description string, durationMinutes int) (*models.MerchantEvent, error) {
        if database.DB == nil {
                return storage.Memory.CreateMerchant("random_shop", durationMinutes)
        }

        // End any existing active merchant events
        _, err := database.DB.Exec("UPDATE merchant_events SET is_active = false WHERE is_active = true")
        if err != nil {
                return nil, fmt.Errorf("failed to deactivate existing events: %v", err)
        }

        // Create the merchant event
        startTime := time.Now()
        endTime := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
        query := `
                INSERT INTO merchant_events (event_type, available_items, start_time, end_time, is_active)
                VALUES (?, ?, ?, ?, true)`

        // Simplified available_items as JSON string
        availableItems := `[]` // Empty for now
        result, err := database.DB.Exec(query, "random_shop", availableItems, startTime, endTime)
        if err != nil {
                return nil, fmt.Errorf("failed to create merchant event: %v", err)
        }

        eventID, err := result.LastInsertId()
        if err != nil {
                return nil, fmt.Errorf("failed to get merchant event ID: %v", err)
        }

        // Get random special items for the event
        itemService := NewItemService()
        randomItems, err := itemService.GetRandomItems(3, true) // 3 special items
        if err != nil {
                return nil, fmt.Errorf("failed to get random items: %v", err)
        }

        // Add items to the merchant event
        for _, item := range randomItems {
                price := item.Value * 2 // Double the item value as channel points price
                stock := rand.Intn(3) + 1 // 1-3 stock

                itemQuery := `
                        INSERT INTO merchant_event_items (merchant_event_id, item_id, price_channel_points, stock, purchased)
                        VALUES (?, ?, ?, ?, 0)`

                _, err = database.DB.Exec(itemQuery, eventID, item.ID, price, stock)
                if err != nil {
                        return nil, fmt.Errorf("failed to add item to merchant event: %v", err)
                }
        }

        return ms.GetMerchantEventByID(int(eventID))
}

// PurchaseItem handles a purchase from a merchant event
func (ms *MerchantService) PurchaseItem(characterID, merchantEventItemID int) error {
        if database.DB == nil {
                // Use item ID directly for simplicity in memory mode
                itemID := merchantEventItemID // Simplified mapping
                price := 100 // Fixed price for testing
                return storage.Memory.PurchaseItem(characterID, itemID, price)
        }

        // Get the merchant event item
        var itemID, price, stock, purchased int
        query := `
                SELECT item_id, price_channel_points, stock, purchased
                FROM merchant_event_items
                WHERE id = ?`

        err := database.DB.QueryRow(query, merchantEventItemID).Scan(&itemID, &price, &stock, &purchased)
        if err != nil {
                return fmt.Errorf("failed to get merchant event item: %v", err)
        }

        // Check if item is still in stock
        if purchased >= stock {
                return fmt.Errorf("item is out of stock")
        }

        // Get character to check channel points
        charService := NewCharacterService()
        character, err := charService.GetCharacterByID(characterID)
        if err != nil {
                return fmt.Errorf("failed to get character: %v", err)
        }

        if character == nil {
                return fmt.Errorf("character not found")
        }

        // Check if character has enough channel points (simplified logic)
        if character.ChannelPointsSpent+price > character.ChannelPointsSpent+1000 { // Simplified check
                return fmt.Errorf("insufficient channel points")
        }

        // Update purchased count
        _, err = database.DB.Exec("UPDATE merchant_event_items SET purchased = purchased + 1 WHERE id = ?", merchantEventItemID)
        if err != nil {
                return fmt.Errorf("failed to update purchase count: %v", err)
        }

        // Add item to character's inventory
        itemService := NewItemService()
        err = itemService.AddItemToCharacter(characterID, itemID, 1)
        if err != nil {
                return fmt.Errorf("failed to add item to character: %v", err)
        }

        // Update character's channel points spent
        character.ChannelPointsSpent += price
        err = charService.UpdateCharacter(character)
        if err != nil {
                return fmt.Errorf("failed to update character: %v", err)
        }

        return nil
}

// GetMerchantEventByID retrieves a merchant event by ID
func (ms *MerchantService) GetMerchantEventByID(id int) (*models.MerchantEvent, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }

        // This is a simplified version - in reality you'd do a proper join query
        query := `
                SELECT id, event_type, available_items, start_time, end_time, is_active
                FROM merchant_events 
                WHERE id = ?`

        event := &models.MerchantEvent{}
        err := database.DB.QueryRow(query, id).Scan(
                &event.ID, &event.EventType, &event.AvailableItems,
                &event.StartTime, &event.EndTime, &event.IsActive,
        )

        if err != nil {
                return nil, fmt.Errorf("failed to get merchant event: %v", err)
        }

        return event, nil
}