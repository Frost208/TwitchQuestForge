package models

import (
	"encoding/json"
	"time"
)

// GameEventType represents different types of game events
type GameEventType string

const (
	EventTypeCombat        GameEventType = "combat"
	EventTypeMerchant      GameEventType = "merchant"
	EventTypeLevelUp       GameEventType = "level_up"
	EventTypeItemAcquired  GameEventType = "item_acquired"
	EventTypeQuestCompleted GameEventType = "quest_completed"
)

// GameEvent represents an event that can trigger OBS animations
type GameEvent struct {
	ID           int           `json:"id" db:"id"`
	EventType    GameEventType `json:"event_type" db:"event_type"`
	CharacterID  *int          `json:"character_id,omitempty" db:"character_id"`
	EventData    json.RawMessage `json:"event_data" db:"event_data"`
	OBSTriggered bool          `json:"obs_triggered" db:"obs_triggered"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	
	// Populated fields
	Character *Character `json:"character,omitempty"`
}

// MerchantEvent represents a merchant appearance event
type MerchantEvent struct {
	ID             int             `json:"id" db:"id"`
	EventType      string          `json:"event_type" db:"event_type"` // 'random_shop', 'special_trader'
	AvailableItems json.RawMessage `json:"available_items" db:"available_items"` // Array of item IDs
	StartTime      time.Time       `json:"start_time" db:"start_time"`
	EndTime        *time.Time      `json:"end_time,omitempty" db:"end_time"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	
	// Populated fields
	Items []Item `json:"items,omitempty"`
}

// Quest represents a quest that characters can complete
type Quest struct {
	ID               int             `json:"id" db:"id"`
	Name             string          `json:"name" db:"name"`
	Description      string          `json:"description" db:"description"`
	QuestType        string          `json:"quest_type" db:"quest_type"` // 'daily', 'weekly', 'special'
	Requirements     json.RawMessage `json:"requirements" db:"requirements"`
	Rewards          json.RawMessage `json:"rewards" db:"rewards"`
	ChannelPointCost int             `json:"channel_point_cost" db:"channel_point_cost"`
	IsActive         bool            `json:"is_active" db:"is_active"`
}

// CharacterQuest represents a character's progress on a quest
type CharacterQuest struct {
	ID          int             `json:"id" db:"id"`
	CharacterID int             `json:"character_id" db:"character_id"`
	QuestID     int             `json:"quest_id" db:"quest_id"`
	Progress    json.RawMessage `json:"progress" db:"progress"`
	Completed   bool            `json:"completed" db:"completed"`
	CompletedAt *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	StartedAt   time.Time       `json:"started_at" db:"started_at"`
	
	// Populated fields
	Quest     *Quest     `json:"quest,omitempty"`
	Character *Character `json:"character,omitempty"`
}

// OBSEventData represents data sent to OBS for animations
type OBSEventData struct {
	EventType   GameEventType          `json:"event_type"`
	CharacterName string              `json:"character_name"`
	Message     string                `json:"message"`
	Data        map[string]interface{} `json:"data"`
}

// CombatEventData represents data for combat events
type CombatEventData struct {
	AttackerName  string `json:"attacker_name"`
	DefenderName  string `json:"defender_name"`
	WinnerName    string `json:"winner_name"`
	AttackerPower int    `json:"attacker_power"`
	DefenderPower int    `json:"defender_power"`
	CombatLog     string `json:"combat_log"`
}

// MerchantEventData represents data for merchant events
type MerchantEventData struct {
	MerchantType  string   `json:"merchant_type"`
	AvailableItems []string `json:"available_items"`
	Duration      int      `json:"duration_minutes"`
}

// LevelUpEventData represents data for level up events
type LevelUpEventData struct {
	CharacterName string `json:"character_name"`
	NewLevel      int    `json:"new_level"`
	NewStats      Stats  `json:"new_stats"`
}

// ItemAcquiredEventData represents data for item acquisition events
type ItemAcquiredEventData struct {
	CharacterName string     `json:"character_name"`
	ItemName      string     `json:"item_name"`
	ItemRarity    ItemRarity `json:"item_rarity"`
	Method        string     `json:"method"` // 'purchase', 'quest_reward', 'combat_reward'
}

// CreateGameEvent creates a new game event
func CreateGameEvent(eventType GameEventType, characterID *int, data interface{}) (*GameEvent, error) {
	eventData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	return &GameEvent{
		EventType:    eventType,
		CharacterID:  characterID,
		EventData:    eventData,
		OBSTriggered: false,
		CreatedAt:    time.Now(),
	}, nil
}

// CreateCombatEvent creates a combat event for OBS
func CreateCombatEvent(combat *CombatResult) (*GameEvent, error) {
	data := CombatEventData{
		AttackerName:  combat.Winner.Username,
		DefenderName:  combat.Loser.Username,
		WinnerName:    combat.Winner.Username,
		AttackerPower: combat.AttackerPower,
		DefenderPower: combat.DefenderPower,
		CombatLog:     combat.CombatLog,
	}
	
	return CreateGameEvent(EventTypeCombat, &combat.Winner.ID, data)
}

// CreateLevelUpEvent creates a level up event for OBS
func CreateLevelUpEvent(character *Character) (*GameEvent, error) {
	data := LevelUpEventData{
		CharacterName: character.Username,
		NewLevel:      character.Level,
		NewStats:      character.CalculateTotalStats(),
	}
	
	return CreateGameEvent(EventTypeLevelUp, &character.ID, data)
}

// CreateItemAcquiredEvent creates an item acquisition event for OBS
func CreateItemAcquiredEvent(character *Character, item *Item, method string) (*GameEvent, error) {
	data := ItemAcquiredEventData{
		CharacterName: character.Username,
		ItemName:      item.Name,
		ItemRarity:    item.Rarity,
		Method:        method,
	}
	
	return CreateGameEvent(EventTypeItemAcquired, &character.ID, data)
}