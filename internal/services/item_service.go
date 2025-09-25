package services

import (
        "database/sql"
        "fmt"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/storage"
)

// ItemService handles item-related operations
type ItemService struct{}

// NewItemService creates a new item service
func NewItemService() *ItemService {
        return &ItemService{}
}

// GetItemByID retrieves an item by ID
func (is *ItemService) GetItemByID(id int) (*models.Item, error) {
        if database.DB == nil {
                return storage.Memory.GetItemByID(id)
        }
        
        query := `
                SELECT id, name, type, rarity, strength_bonus, agility_bonus, 
                        vitality_bonus, intelligence_bonus, special_effect, value, is_special, created_at
                FROM items WHERE id = ?`
        
        item := &models.Item{}
        err := database.DB.QueryRow(query, id).Scan(
                &item.ID, &item.Name, &item.Type, &item.Rarity,
                &item.StrengthBonus, &item.AgilityBonus, &item.VitalityBonus, &item.IntelligenceBonus,
                &item.SpecialEffect, &item.Value, &item.IsSpecial, &item.CreatedAt,
        )
        
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil
                }
                return nil, fmt.Errorf("failed to get item: %v", err)
        }
        
        return item, nil
}

// GetItemsByType retrieves items by type with pagination
func (is *ItemService) GetItemsByType(itemType models.ItemType, limit, offset int) ([]models.Item, error) {
        if database.DB == nil {
                return storage.Memory.GetItemsByType(itemType, limit, offset)
        }
        
        query := `
                SELECT id, name, type, rarity, strength_bonus, agility_bonus, 
                        vitality_bonus, intelligence_bonus, special_effect, value, is_special, created_at
                FROM items WHERE type = ? ORDER BY rarity DESC, value DESC LIMIT ? OFFSET ?`
        
        rows, err := database.DB.Query(query, itemType, limit, offset)
        if err != nil {
                return nil, fmt.Errorf("failed to get items: %v", err)
        }
        defer rows.Close()
        
        var items []models.Item
        for rows.Next() {
                var item models.Item
                err := rows.Scan(
                        &item.ID, &item.Name, &item.Type, &item.Rarity,
                        &item.StrengthBonus, &item.AgilityBonus, &item.VitalityBonus, &item.IntelligenceBonus,
                        &item.SpecialEffect, &item.Value, &item.IsSpecial, &item.CreatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan item: %v", err)
                }
                items = append(items, item)
        }
        
        return items, nil
}

// GetRandomItems retrieves random items for merchant events
func (is *ItemService) GetRandomItems(count int, isSpecial bool) ([]models.Item, error) {
        if database.DB == nil {
                return storage.Memory.GetRandomItems(count, isSpecial)
        }
        
        query := `
                SELECT id, name, type, rarity, strength_bonus, agility_bonus, 
                        vitality_bonus, intelligence_bonus, special_effect, value, is_special, created_at
                FROM items WHERE is_special = ? ORDER BY RAND() LIMIT ?`
        
        rows, err := database.DB.Query(query, isSpecial, count)
        if err != nil {
                return nil, fmt.Errorf("failed to get random items: %v", err)
        }
        defer rows.Close()
        
        var items []models.Item
        for rows.Next() {
                var item models.Item
                err := rows.Scan(
                        &item.ID, &item.Name, &item.Type, &item.Rarity,
                        &item.StrengthBonus, &item.AgilityBonus, &item.VitalityBonus, &item.IntelligenceBonus,
                        &item.SpecialEffect, &item.Value, &item.IsSpecial, &item.CreatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan item: %v", err)
                }
                items = append(items, item)
        }
        
        return items, nil
}

// CharacterOwnsItem checks if a character owns a specific item
func (is *ItemService) CharacterOwnsItem(characterID, itemID int) (bool, error) {
        if database.DB == nil {
                // Simplified: assume character owns any item for testing
                return true, nil
        }
        
        query := `SELECT COUNT(*) FROM character_items WHERE character_id = ? AND item_id = ?`
        
        var count int
        err := database.DB.QueryRow(query, characterID, itemID).Scan(&count)
        if err != nil {
                return false, fmt.Errorf("failed to check item ownership: %v", err)
        }
        
        return count > 0, nil
}

// AddItemToCharacter adds an item to a character's inventory
func (is *ItemService) AddItemToCharacter(characterID, itemID, quantity int) error {
        if database.DB == nil {
                return fmt.Errorf("database connection not available")
        }
        
        // Check if character already has this item
        query := `SELECT quantity FROM character_items WHERE character_id = ? AND item_id = ?`
        
        var currentQuantity int
        err := database.DB.QueryRow(query, characterID, itemID).Scan(&currentQuantity)
        
        if err == sql.ErrNoRows {
                // Character doesn't have this item, create new entry
                insertQuery := `INSERT INTO character_items (character_id, item_id, quantity) VALUES (?, ?, ?)`
                _, err = database.DB.Exec(insertQuery, characterID, itemID, quantity)
                if err != nil {
                        return fmt.Errorf("failed to add item to character: %v", err)
                }
        } else if err != nil {
                return fmt.Errorf("failed to check existing item: %v", err)
        } else {
                // Character already has this item, update quantity
                updateQuery := `UPDATE character_items SET quantity = quantity + ? WHERE character_id = ? AND item_id = ?`
                _, err = database.DB.Exec(updateQuery, quantity, characterID, itemID)
                if err != nil {
                        return fmt.Errorf("failed to update item quantity: %v", err)
                }
        }
        
        return nil
}

// GetCharacterItems retrieves all items owned by a character
func (is *ItemService) GetCharacterItems(characterID int) ([]models.CharacterItem, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }
        
        query := `
                SELECT ci.id, ci.character_id, ci.item_id, ci.quantity, ci.acquired_at,
                        i.id, i.name, i.type, i.rarity, i.strength_bonus, i.agility_bonus, 
                        i.vitality_bonus, i.intelligence_bonus, i.special_effect, i.value, i.is_special, i.created_at
                FROM character_items ci
                JOIN items i ON ci.item_id = i.id
                WHERE ci.character_id = ?
                ORDER BY i.rarity DESC, i.value DESC`
        
        rows, err := database.DB.Query(query, characterID)
        if err != nil {
                return nil, fmt.Errorf("failed to get character items: %v", err)
        }
        defer rows.Close()
        
        var characterItems []models.CharacterItem
        for rows.Next() {
                var ci models.CharacterItem
                var item models.Item
                
                err := rows.Scan(
                        &ci.ID, &ci.CharacterID, &ci.ItemID, &ci.Quantity, &ci.AcquiredAt,
                        &item.ID, &item.Name, &item.Type, &item.Rarity,
                        &item.StrengthBonus, &item.AgilityBonus, &item.VitalityBonus, &item.IntelligenceBonus,
                        &item.SpecialEffect, &item.Value, &item.IsSpecial, &item.CreatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan character item: %v", err)
                }
                
                ci.Item = &item
                characterItems = append(characterItems, ci)
        }
        
        return characterItems, nil
}