package storage

import (
        "fmt"
        "sync"
        "time"
        "twitch-rpg/internal/models"
)

// MemoryStorage provides in-memory storage for testing when database is unavailable
type MemoryStorage struct {
        characters     map[int]*models.Character
        items          map[int]*models.Item
        combatLogs     []models.CombatLog
        events         []models.Event
        merchants      []models.MerchantEvent
        activeMerchant *models.MerchantEvent
        inventories    map[int][]models.Item  // characterID -> items
        
        nextCharacterID int
        nextCombatLogID int
        nextEventID     int
        nextMerchantID  int
        
        mutex sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
        ms := &MemoryStorage{
                characters:      make(map[int]*models.Character),
                items:           make(map[int]*models.Item),
                combatLogs:      []models.CombatLog{},
                events:          []models.Event{},
                merchants:       []models.MerchantEvent{},
                activeMerchant:   nil,
                inventories:      make(map[int][]models.Item),
                nextCharacterID: 1,
                nextCombatLogID: 1,
                nextEventID:     1,
                nextMerchantID:  1,
        }
        
        // Initialize with sample data
        ms.initializeSampleData()
        return ms
}

// Global instance
var Memory *MemoryStorage

func init() {
        Memory = NewMemoryStorage()
}

// Character operations
func (ms *MemoryStorage) CreateCharacter(username string, twitchUserID *string) (*models.Character, error) {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        // Check if username exists
        for _, char := range ms.characters {
                if char.Username == username {
                        return nil, fmt.Errorf("character with username '%s' already exists", username)
                }
        }
        
        char := &models.Character{
                ID:                 ms.nextCharacterID,
                Username:           username,
                TwitchUserID:       twitchUserID,
                Level:              1,
                Experience:         0,
                ChannelPointsSpent: 0,
                Strength:           10,
                Agility:            10,
                Vitality:           10,
                Intelligence:       10,
                CreatedAt:          time.Now(),
                UpdatedAt:          time.Now(),
        }
        
        char.CombatPower = char.CalculateCombatPower()
        
        ms.characters[ms.nextCharacterID] = char
        ms.nextCharacterID++
        
        return char, nil
}

func (ms *MemoryStorage) GetCharacterByID(id int) (*models.Character, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        char, exists := ms.characters[id]
        if !exists {
                return nil, nil
        }
        
        // Create a copy to avoid race conditions
        result := *char
        result.CombatPower = result.CalculateCombatPower()
        
        return &result, nil
}

func (ms *MemoryStorage) GetCharacterByUsername(username string) (*models.Character, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        for _, char := range ms.characters {
                if char.Username == username {
                        result := *char
                        result.CombatPower = result.CalculateCombatPower()
                        return &result, nil
                }
        }
        
        return nil, nil
}

func (ms *MemoryStorage) UpdateCharacter(char *models.Character) error {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        if _, exists := ms.characters[char.ID]; !exists {
                return fmt.Errorf("character not found")
        }
        
        char.UpdatedAt = time.Now()
        ms.characters[char.ID] = char
        
        return nil
}

func (ms *MemoryStorage) GetAllCharacters() ([]models.Character, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        var characters []models.Character
        for _, char := range ms.characters {
                result := *char
                result.CombatPower = result.CalculateCombatPower()
                characters = append(characters, result)
        }
        
        return characters, nil
}

// Item operations
func (ms *MemoryStorage) GetItemByID(id int) (*models.Item, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        item, exists := ms.items[id]
        if !exists {
                return nil, nil
        }
        
        result := *item
        return &result, nil
}

func (ms *MemoryStorage) GetRandomItems(count int, isSpecial bool) ([]models.Item, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        var items []models.Item
        added := 0
        
        for _, item := range ms.items {
                if item.IsSpecial == isSpecial && added < count {
                        items = append(items, *item)
                        added++
                }
        }
        
        return items, nil
}

func (ms *MemoryStorage) GetItemsByType(itemType models.ItemType, limit, offset int) ([]models.Item, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        var items []models.Item
        skipped := 0
        added := 0
        
        for _, item := range ms.items {
                if item.Type == itemType {
                        if skipped < offset {
                                skipped++
                                continue
                        }
                        if added >= limit {
                                break
                        }
                        items = append(items, *item)
                        added++
                }
        }
        
        return items, nil
}

// Combat operations
func (ms *MemoryStorage) AddCombatLog(log models.CombatLog) error {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        log.ID = ms.nextCombatLogID
        log.CreatedAt = time.Now()
        ms.combatLogs = append(ms.combatLogs, log)
        ms.nextCombatLogID++
        
        return nil
}

func (ms *MemoryStorage) GetCombatHistory(limit int) ([]models.CombatLog, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        // Return most recent logs
        start := len(ms.combatLogs) - limit
        if start < 0 {
                start = 0
        }
        
        var result []models.CombatLog
        for i := len(ms.combatLogs) - 1; i >= start; i-- {
                result = append(result, ms.combatLogs[i])
        }
        
        return result, nil
}

// Event operations
func (ms *MemoryStorage) GetLatestEvents(limit int) ([]models.Event, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        // Return most recent events
        start := len(ms.events) - limit
        if start < 0 {
                start = 0
        }
        
        var result []models.Event
        for i := len(ms.events) - 1; i >= start; i-- {
                result = append(result, ms.events[i])
        }
        
        return result, nil
}

func (ms *MemoryStorage) MarkEventTriggered(eventID int) error {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        for i := range ms.events {
                if ms.events[i].ID == eventID {
                        ms.events[i].IsTriggered = true
                        return nil
                }
        }
        
        return fmt.Errorf("event not found")
}

// Initialize sample data for testing
func (ms *MemoryStorage) initializeSampleData() {
        // Sample items
        ms.items[1] = &models.Item{
                ID:               1,
                Name:             "Iron Boots",
                Type:             models.ItemTypeBoots,
                Rarity:           models.RarityCommon,
                StrengthBonus:    2,
                AgilityBonus:     1,
                VitalityBonus:    1,
                IntelligenceBonus: 0,
                Value:            50,
                IsSpecial:        false,
                CreatedAt:        time.Now(),
        }
        
        ms.items[2] = &models.Item{
                ID:               2,
                Name:             "Mystic Sword",
                Type:             models.ItemTypeRing, // Using ring as example
                Rarity:           models.RarityRare,
                StrengthBonus:    5,
                AgilityBonus:     3,
                VitalityBonus:    2,
                IntelligenceBonus: 4,
                SpecialEffect:    stringPtr("Increases magic damage"),
                Value:            200,
                IsSpecial:        true,
                CreatedAt:        time.Now(),
        }
        
        ms.items[3] = &models.Item{
                ID:               3,
                Name:             "Leather Pants",
                Type:             models.ItemTypePants,
                Rarity:           models.RarityCommon,
                StrengthBonus:    1,
                AgilityBonus:     3,
                VitalityBonus:    2,
                IntelligenceBonus: 0,
                Value:            40,
                IsSpecial:        false,
                CreatedAt:        time.Now(),
        }
        
        // Sample events
        ms.events = append(ms.events, models.Event{
                ID:          1,
                Type:        "combat",
                Title:       "Epic Battle",
                Description: "A fierce battle took place in the arena",
                Data:        `{"winner": "TestUser", "experience": 100}`,
                IsTriggered: false,
                CreatedAt:   time.Now(),
        })
}

// Merchant operations
func (ms *MemoryStorage) GetCurrentMerchant() (*models.MerchantEvent, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        if ms.activeMerchant != nil && ms.activeMerchant.IsActive {
                // Check if still active
                if ms.activeMerchant.EndTime == nil || ms.activeMerchant.EndTime.After(time.Now()) {
                        result := *ms.activeMerchant
                        return &result, nil
                }
                // Expired
                ms.activeMerchant = nil
        }
        
        return nil, nil
}

func (ms *MemoryStorage) CreateMerchant(eventType string, durationMinutes int) (*models.MerchantEvent, error) {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        // Deactivate current merchant
        if ms.activeMerchant != nil {
                ms.activeMerchant.IsActive = false
        }
        
        // Create new merchant event
        endTime := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
        merchant := &models.MerchantEvent{
                ID:             ms.nextMerchantID,
                EventType:      eventType,
                AvailableItems: []byte(`[1,2,3]`), // Sample items
                StartTime:      time.Now(),
                EndTime:        &endTime,
                IsActive:       true,
        }
        
        ms.activeMerchant = merchant
        ms.merchants = append(ms.merchants, *merchant)
        ms.nextMerchantID++
        
        return merchant, nil
}

func (ms *MemoryStorage) PurchaseItem(characterID, itemID int, price int) error {
        ms.mutex.Lock()
        defer ms.mutex.Unlock()
        
        // Get character
        char, exists := ms.characters[characterID]
        if !exists {
                return fmt.Errorf("character not found")
        }
        
        // Get item
        item, exists := ms.items[itemID]
        if !exists {
                return fmt.Errorf("item not found")
        }
        
        // Update character points
        char.ChannelPointsSpent += price
        
        // Add item to character's inventory
        if ms.inventories[characterID] == nil {
                ms.inventories[characterID] = []models.Item{}
        }
        ms.inventories[characterID] = append(ms.inventories[characterID], *item)
        
        return nil
}

func (ms *MemoryStorage) GetCharacterInventory(characterID int) ([]models.Item, error) {
        ms.mutex.RLock()
        defer ms.mutex.RUnlock()
        
        if inventory, exists := ms.inventories[characterID]; exists {
                return inventory, nil
        }
        
        return []models.Item{}, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
        return &s
}