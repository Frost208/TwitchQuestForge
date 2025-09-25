package models

import (
        "fmt"
        "strings"
        "time"
)

// ItemType represents the type of equipment item
type ItemType string

const (
        ItemTypeBoots   ItemType = "boots"
        ItemTypePants   ItemType = "pants"
        ItemTypeArmor   ItemType = "armor"
        ItemTypeHelmet  ItemType = "helmet"
        ItemTypeRing    ItemType = "ring"
        ItemTypeChain   ItemType = "chain"
)

// ItemRarity represents the rarity level of an item
type ItemRarity string

const (
        RarityCommon    ItemRarity = "common"
        RarityRare      ItemRarity = "rare"
        RarityEpic      ItemRarity = "epic"
        RarityLegendary ItemRarity = "legendary"
)

// Item represents an equipment item in the game
type Item struct {
        ID               int         `json:"id" db:"id"`
        Name             string      `json:"name" db:"name"`
        Type             ItemType    `json:"type" db:"type"`
        Rarity           ItemRarity  `json:"rarity" db:"rarity"`
        StrengthBonus    int         `json:"strength_bonus" db:"strength_bonus"`
        AgilityBonus     int         `json:"agility_bonus" db:"agility_bonus"`
        VitalityBonus    int         `json:"vitality_bonus" db:"vitality_bonus"`
        IntelligenceBonus int        `json:"intelligence_bonus" db:"intelligence_bonus"`
        SpecialEffect    *string     `json:"special_effect,omitempty" db:"special_effect"`
        Value            int         `json:"value" db:"value"` // Channel points value
        IsSpecial        bool        `json:"is_special" db:"is_special"` // For merchant items
        CreatedAt        time.Time   `json:"created_at" db:"created_at"`
}

// CharacterItem represents an item owned by a character
type CharacterItem struct {
        ID          int       `json:"id" db:"id"`
        CharacterID int       `json:"character_id" db:"character_id"`
        ItemID      int       `json:"item_id" db:"item_id"`
        Quantity    int       `json:"quantity" db:"quantity"`
        AcquiredAt  time.Time `json:"acquired_at" db:"acquired_at"`
        
        // Populated fields
        Item *Item `json:"item,omitempty"`
}

// ItemFilter represents filters for querying items
type ItemFilter struct {
        Type      *ItemType    `json:"type,omitempty"`
        Rarity    *ItemRarity  `json:"rarity,omitempty"`
        MinValue  *int         `json:"min_value,omitempty"`
        MaxValue  *int         `json:"max_value,omitempty"`
        IsSpecial *bool        `json:"is_special,omitempty"`
        Limit     int          `json:"limit"`
        Offset    int          `json:"offset"`
}

// EquipItemRequest represents a request to equip an item
type EquipItemRequest struct {
        ItemID int `json:"item_id" binding:"required"`
}

// GetTotalStatBonus calculates the total stat bonus of an item
func (i *Item) GetTotalStatBonus() int {
        return i.StrengthBonus + i.AgilityBonus + i.VitalityBonus + i.IntelligenceBonus
}

// GetRarityMultiplier returns a multiplier based on item rarity
func (i *Item) GetRarityMultiplier() float64 {
        switch i.Rarity {
        case RarityCommon:
                return 1.0
        case RarityRare:
                return 1.5
        case RarityEpic:
                return 2.0
        case RarityLegendary:
                return 3.0
        default:
                return 1.0
        }
}

// GetEffectiveValue calculates the effective channel point value including rarity
func (i *Item) GetEffectiveValue() int {
        return int(float64(i.Value) * i.GetRarityMultiplier())
}

// CanEquipToSlot checks if this item can be equipped to a specific slot
func (i *Item) CanEquipToSlot(slotType ItemType) bool {
        return i.Type == slotType
}

// IsUpgradeFor checks if this item is an upgrade compared to another item
func (i *Item) IsUpgradeFor(other *Item) bool {
        if other == nil {
                return true // Any item is better than no item
        }
        
        if i.Type != other.Type {
                return false // Different types can't be compared
        }
        
        // Compare total stat bonus first
        thisTotalBonus := i.GetTotalStatBonus()
        otherTotalBonus := other.GetTotalStatBonus()
        
        if thisTotalBonus != otherTotalBonus {
                return thisTotalBonus > otherTotalBonus
        }
        
        // If total bonus is equal, compare rarity
        rarityOrder := map[ItemRarity]int{
                RarityCommon:    1,
                RarityRare:      2,
                RarityEpic:      3,
                RarityLegendary: 4,
        }
        
        return rarityOrder[i.Rarity] > rarityOrder[other.Rarity]
}

// GetDescription generates a description of the item
func (i *Item) GetDescription() string {
        description := i.Name
        
        if i.SpecialEffect != nil && *i.SpecialEffect != "" {
                description += " - " + *i.SpecialEffect
        }
        
        // Add stat bonuses
        bonuses := []string{}
        if i.StrengthBonus > 0 {
                bonuses = append(bonuses, fmt.Sprintf("STR +%d", i.StrengthBonus))
        }
        if i.AgilityBonus > 0 {
                bonuses = append(bonuses, fmt.Sprintf("AGI +%d", i.AgilityBonus))
        }
        if i.VitalityBonus > 0 {
                bonuses = append(bonuses, fmt.Sprintf("VIT +%d", i.VitalityBonus))
        }
        if i.IntelligenceBonus > 0 {
                bonuses = append(bonuses, fmt.Sprintf("INT +%d", i.IntelligenceBonus))
        }
        
        if len(bonuses) > 0 {
                description += " (" + strings.Join(bonuses, ", ") + ")"
        }
        
        return description
}

// ValidateItemType checks if a string is a valid ItemType
func ValidateItemType(itemType string) bool {
        switch ItemType(itemType) {
        case ItemTypeBoots, ItemTypePants, ItemTypeArmor, ItemTypeHelmet, ItemTypeRing, ItemTypeChain:
                return true
        default:
                return false
        }
}

// ValidateItemRarity checks if a string is a valid ItemRarity
func ValidateItemRarity(rarity string) bool {
        switch ItemRarity(rarity) {
        case RarityCommon, RarityRare, RarityEpic, RarityLegendary:
                return true
        default:
                return false
        }
}