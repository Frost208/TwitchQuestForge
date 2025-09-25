package services

import (
        "database/sql"
        "fmt"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/models"
)

// CharacterService handles character-related operations
type CharacterService struct{}

// NewCharacterService creates a new character service
func NewCharacterService() *CharacterService {
        return &CharacterService{}
}

// CreateCharacter creates a new character for a user
func (cs *CharacterService) CreateCharacter(username string, twitchUserID *string) (*models.Character, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }
        
        // Check if character already exists
        existing, _ := cs.GetCharacterByUsername(username)
        if existing != nil {
                return nil, fmt.Errorf("character with username '%s' already exists", username)
        }
        
        query := `
                INSERT INTO characters (username, twitch_user_id, level, experience, channel_points_spent, 
                        strength, agility, vitality, intelligence) 
                VALUES (?, ?, 1, 0, 0, 10, 10, 10, 10)`
        
        result, err := database.DB.Exec(query, username, twitchUserID)
        if err != nil {
                return nil, fmt.Errorf("failed to create character: %v", err)
        }
        
        id, err := result.LastInsertId()
        if err != nil {
                return nil, fmt.Errorf("failed to get character ID: %v", err)
        }
        
        return cs.GetCharacterByID(int(id))
}

// GetCharacterByID retrieves a character by ID
func (cs *CharacterService) GetCharacterByID(id int) (*models.Character, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }
        query := `
                SELECT id, username, twitch_user_id, level, experience, channel_points_spent,
                        strength, agility, vitality, intelligence,
                        boots_id, pants_id, armor_id, helmet_id, ring_id, chain_id,
                        created_at, updated_at
                FROM characters WHERE id = ?`
        
        character := &models.Character{}
        err := database.DB.QueryRow(query, id).Scan(
                &character.ID, &character.Username, &character.TwitchUserID,
                &character.Level, &character.Experience, &character.ChannelPointsSpent,
                &character.Strength, &character.Agility, &character.Vitality, &character.Intelligence,
                &character.BootsID, &character.PantsID, &character.ArmorID,
                &character.HelmetID, &character.RingID, &character.ChainID,
                &character.CreatedAt, &character.UpdatedAt,
        )
        
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil
                }
                return nil, fmt.Errorf("failed to get character: %v", err)
        }
        
        // Load equipment
        if err := cs.loadCharacterEquipment(character); err != nil {
                return nil, fmt.Errorf("failed to load equipment: %v", err)
        }
        
        // Calculate derived stats
        totalStats := character.CalculateTotalStats()
        character.TotalStats = &totalStats
        character.CombatPower = character.CalculateCombatPower()
        
        return character, nil
}

// GetCharacterByUsername retrieves a character by username
func (cs *CharacterService) GetCharacterByUsername(username string) (*models.Character, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }
        query := `
                SELECT id, username, twitch_user_id, level, experience, channel_points_spent,
                        strength, agility, vitality, intelligence,
                        boots_id, pants_id, armor_id, helmet_id, ring_id, chain_id,
                        created_at, updated_at
                FROM characters WHERE username = ?`
        
        character := &models.Character{}
        err := database.DB.QueryRow(query, username).Scan(
                &character.ID, &character.Username, &character.TwitchUserID,
                &character.Level, &character.Experience, &character.ChannelPointsSpent,
                &character.Strength, &character.Agility, &character.Vitality, &character.Intelligence,
                &character.BootsID, &character.PantsID, &character.ArmorID,
                &character.HelmetID, &character.RingID, &character.ChainID,
                &character.CreatedAt, &character.UpdatedAt,
        )
        
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil
                }
                return nil, fmt.Errorf("failed to get character: %v", err)
        }
        
        // Load equipment
        if err := cs.loadCharacterEquipment(character); err != nil {
                return nil, fmt.Errorf("failed to load equipment: %v", err)
        }
        
        // Calculate derived stats
        totalStats := character.CalculateTotalStats()
        character.TotalStats = &totalStats
        character.CombatPower = character.CalculateCombatPower()
        
        return character, nil
}

// UpdateCharacter updates character information
func (cs *CharacterService) UpdateCharacter(character *models.Character) error {
        query := `
                UPDATE characters SET 
                        level = ?, experience = ?, channel_points_spent = ?,
                        strength = ?, agility = ?, vitality = ?, intelligence = ?,
                        boots_id = ?, pants_id = ?, armor_id = ?, helmet_id = ?, ring_id = ?, chain_id = ?
                WHERE id = ?`
        
        _, err := database.DB.Exec(query,
                character.Level, character.Experience, character.ChannelPointsSpent,
                character.Strength, character.Agility, character.Vitality, character.Intelligence,
                character.BootsID, character.PantsID, character.ArmorID,
                character.HelmetID, character.RingID, character.ChainID,
                character.ID,
        )
        
        if err != nil {
                return fmt.Errorf("failed to update character: %v", err)
        }
        
        return nil
}

// UpgradeCharacterStat upgrades a character's stat using channel points
func (cs *CharacterService) UpgradeCharacterStat(characterID int, statType string, channelPoints int) (*models.Character, error) {
        character, err := cs.GetCharacterByID(characterID)
        if err != nil {
                return nil, err
        }
        
        if character == nil {
                return nil, fmt.Errorf("character not found")
        }
        
        if !character.UpgradeStat(statType, channelPoints) {
                return nil, fmt.Errorf("invalid stat upgrade parameters")
        }
        
        if err := cs.UpdateCharacter(character); err != nil {
                return nil, err
        }
        
        return cs.GetCharacterByID(characterID)
}

// EquipItem equips an item to a character
func (cs *CharacterService) EquipItem(characterID, itemID int) error {
        // Get the item first to check its type
        itemService := NewItemService()
        item, err := itemService.GetItemByID(itemID)
        if err != nil {
                return err
        }
        
        if item == nil {
                return fmt.Errorf("item not found")
        }
        
        // Get the character
        character, err := cs.GetCharacterByID(characterID)
        if err != nil {
                return err
        }
        
        if character == nil {
                return fmt.Errorf("character not found")
        }
        
        // Check if character owns this item
        owns, err := itemService.CharacterOwnsItem(characterID, itemID)
        if err != nil {
                return err
        }
        
        if !owns {
                return fmt.Errorf("character does not own this item")
        }
        
        // Equip the item based on its type
        switch item.Type {
        case models.ItemTypeBoots:
                character.BootsID = &itemID
        case models.ItemTypePants:
                character.PantsID = &itemID
        case models.ItemTypeArmor:
                character.ArmorID = &itemID
        case models.ItemTypeHelmet:
                character.HelmetID = &itemID
        case models.ItemTypeRing:
                character.RingID = &itemID
        case models.ItemTypeChain:
                character.ChainID = &itemID
        default:
                return fmt.Errorf("invalid item type")
        }
        
        return cs.UpdateCharacter(character)
}

// UnequipItem removes an equipped item from a character
func (cs *CharacterService) UnequipItem(characterID int, slotType models.ItemType) error {
        character, err := cs.GetCharacterByID(characterID)
        if err != nil {
                return err
        }
        
        if character == nil {
                return fmt.Errorf("character not found")
        }
        
        // Unequip the item based on slot type
        switch slotType {
        case models.ItemTypeBoots:
                character.BootsID = nil
        case models.ItemTypePants:
                character.PantsID = nil
        case models.ItemTypeArmor:
                character.ArmorID = nil
        case models.ItemTypeHelmet:
                character.HelmetID = nil
        case models.ItemTypeRing:
                character.RingID = nil
        case models.ItemTypeChain:
                character.ChainID = nil
        default:
                return fmt.Errorf("invalid slot type")
        }
        
        return cs.UpdateCharacter(character)
}

// loadCharacterEquipment loads the equipped items for a character
func (cs *CharacterService) loadCharacterEquipment(character *models.Character) error {
        itemService := NewItemService()
        equipment := &models.Equipment{}
        
        var err error
        
        // Load each equipped item
        if character.BootsID != nil {
                equipment.Boots, err = itemService.GetItemByID(*character.BootsID)
                if err != nil {
                        return err
                }
        }
        
        if character.PantsID != nil {
                equipment.Pants, err = itemService.GetItemByID(*character.PantsID)
                if err != nil {
                        return err
                }
        }
        
        if character.ArmorID != nil {
                equipment.Armor, err = itemService.GetItemByID(*character.ArmorID)
                if err != nil {
                        return err
                }
        }
        
        if character.HelmetID != nil {
                equipment.Helmet, err = itemService.GetItemByID(*character.HelmetID)
                if err != nil {
                        return err
                }
        }
        
        if character.RingID != nil {
                equipment.Ring, err = itemService.GetItemByID(*character.RingID)
                if err != nil {
                        return err
                }
        }
        
        if character.ChainID != nil {
                equipment.Chain, err = itemService.GetItemByID(*character.ChainID)
                if err != nil {
                        return err
                }
        }
        
        character.Equipment = equipment
        return nil
}

// GetAllCharacters retrieves all characters with basic info
func (cs *CharacterService) GetAllCharacters() ([]models.Character, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }
        query := `
                SELECT id, username, level, experience, 
                        strength, agility, vitality, intelligence, created_at
                FROM characters 
                ORDER BY level DESC, experience DESC`
        
        rows, err := database.DB.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to get characters: %v", err)
        }
        defer rows.Close()
        
        var characters []models.Character
        for rows.Next() {
                var char models.Character
                err := rows.Scan(
                        &char.ID, &char.Username, &char.Level, &char.Experience,
                        &char.Strength, &char.Agility, &char.Vitality, &char.Intelligence,
                        &char.CreatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan character: %v", err)
                }
                
                // Calculate basic stats
                char.CombatPower = char.CalculateCombatPower()
                characters = append(characters, char)
        }
        
        return characters, nil
}