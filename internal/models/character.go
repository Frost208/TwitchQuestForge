package models

import (
        "time"
)

// Character represents a player character in the Twitch RPG
type Character struct {
        ID                int       `json:"id" db:"id"`
        Username          string    `json:"username" db:"username"`
        TwitchUserID      *string   `json:"twitch_user_id,omitempty" db:"twitch_user_id"`
        Level             int       `json:"level" db:"level"`
        Experience        int       `json:"experience" db:"experience"`
        ChannelPointsSpent int      `json:"channel_points_spent" db:"channel_points_spent"`
        
        // Base stats
        Strength     int `json:"strength" db:"strength"`
        Agility      int `json:"agility" db:"agility"`
        Vitality     int `json:"vitality" db:"vitality"`
        Intelligence int `json:"intelligence" db:"intelligence"`
        
        // Equipment slots (item IDs)
        BootsID   *int `json:"boots_id,omitempty" db:"boots_id"`
        PantsID   *int `json:"pants_id,omitempty" db:"pants_id"`
        ArmorID   *int `json:"armor_id,omitempty" db:"armor_id"`
        HelmetID  *int `json:"helmet_id,omitempty" db:"helmet_id"`
        RingID    *int `json:"ring_id,omitempty" db:"ring_id"`
        ChainID   *int `json:"chain_id,omitempty" db:"chain_id"`
        
        CreatedAt time.Time `json:"created_at" db:"created_at"`
        UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
        
        // Calculated fields (not stored in DB)
        Equipment     *Equipment     `json:"equipment,omitempty"`
        TotalStats    *Stats         `json:"total_stats,omitempty"`
        CombatPower   int            `json:"combat_power,omitempty"`
}

// Stats represents character statistics
type Stats struct {
        Strength     int `json:"strength"`
        Agility      int `json:"agility"`
        Vitality     int `json:"vitality"`
        Intelligence int `json:"intelligence"`
}

// Equipment represents all equipped items for a character
type Equipment struct {
        Boots   *Item `json:"boots,omitempty"`
        Pants   *Item `json:"pants,omitempty"`
        Armor   *Item `json:"armor,omitempty"`
        Helmet  *Item `json:"helmet,omitempty"`
        Ring    *Item `json:"ring,omitempty"`
        Chain   *Item `json:"chain,omitempty"`
}

// CharacterCreateRequest represents the request to create a new character
type CharacterCreateRequest struct {
        Username     string  `json:"username" binding:"required"`
        TwitchUserID *string `json:"twitch_user_id,omitempty"`
}

// CharacterStatsUpgradeRequest represents a request to upgrade character stats
type CharacterStatsUpgradeRequest struct {
        StatType      string `json:"stat_type" binding:"required"` // strength, agility, vitality, intelligence
        ChannelPoints int    `json:"channel_points" binding:"required"`
}

// StatUpgradeRequest is an alias for CharacterStatsUpgradeRequest for backward compatibility
type StatUpgradeRequest = CharacterStatsUpgradeRequest

// CalculateBaseStats calculates the character's base stats without equipment
func (c *Character) CalculateBaseStats() Stats {
        return Stats{
                Strength:     c.Strength,
                Agility:      c.Agility,
                Vitality:     c.Vitality,
                Intelligence: c.Intelligence,
        }
}

// CalculateTotalStats calculates total stats including equipment bonuses
func (c *Character) CalculateTotalStats() Stats {
        baseStats := c.CalculateBaseStats()
        
        if c.Equipment == nil {
                return baseStats
        }
        
        // Add equipment bonuses
        if c.Equipment.Boots != nil {
                baseStats.Strength += c.Equipment.Boots.StrengthBonus
                baseStats.Agility += c.Equipment.Boots.AgilityBonus
                baseStats.Vitality += c.Equipment.Boots.VitalityBonus
                baseStats.Intelligence += c.Equipment.Boots.IntelligenceBonus
        }
        
        if c.Equipment.Pants != nil {
                baseStats.Strength += c.Equipment.Pants.StrengthBonus
                baseStats.Agility += c.Equipment.Pants.AgilityBonus
                baseStats.Vitality += c.Equipment.Pants.VitalityBonus
                baseStats.Intelligence += c.Equipment.Pants.IntelligenceBonus
        }
        
        if c.Equipment.Armor != nil {
                baseStats.Strength += c.Equipment.Armor.StrengthBonus
                baseStats.Agility += c.Equipment.Armor.AgilityBonus
                baseStats.Vitality += c.Equipment.Armor.VitalityBonus
                baseStats.Intelligence += c.Equipment.Armor.IntelligenceBonus
        }
        
        if c.Equipment.Helmet != nil {
                baseStats.Strength += c.Equipment.Helmet.StrengthBonus
                baseStats.Agility += c.Equipment.Helmet.AgilityBonus
                baseStats.Vitality += c.Equipment.Helmet.VitalityBonus
                baseStats.Intelligence += c.Equipment.Helmet.IntelligenceBonus
        }
        
        if c.Equipment.Ring != nil {
                baseStats.Strength += c.Equipment.Ring.StrengthBonus
                baseStats.Agility += c.Equipment.Ring.AgilityBonus
                baseStats.Vitality += c.Equipment.Ring.VitalityBonus
                baseStats.Intelligence += c.Equipment.Ring.IntelligenceBonus
        }
        
        if c.Equipment.Chain != nil {
                baseStats.Strength += c.Equipment.Chain.StrengthBonus
                baseStats.Agility += c.Equipment.Chain.AgilityBonus
                baseStats.Vitality += c.Equipment.Chain.VitalityBonus
                baseStats.Intelligence += c.Equipment.Chain.IntelligenceBonus
        }
        
        return baseStats
}

// CalculateCombatPower calculates the overall combat power of the character
func (c *Character) CalculateCombatPower() int {
        totalStats := c.CalculateTotalStats()
        
        // Combat power formula: weighted sum of all stats + level bonus
        combatPower := (totalStats.Strength * 3) +      // Strength has highest weight for combat
                                        (totalStats.Agility * 2) +        // Agility affects speed and critical hits
                                        (totalStats.Vitality * 2) +       // Vitality affects health and defense
                                        (totalStats.Intelligence * 1) +   // Intelligence affects special abilities
                                        (c.Level * 5)                      // Level provides flat bonus
        
        return combatPower
}

// CanUpgradeStat checks if a character can upgrade a specific stat
func (c *Character) CanUpgradeStat(statType string, channelPoints int) bool {
        // Basic validation - can extend with more complex rules later
        return channelPoints > 0 && (statType == "strength" || statType == "agility" || statType == "vitality" || statType == "intelligence")
}

// UpgradeStat upgrades a character's stat by spending channel points
func (c *Character) UpgradeStat(statType string, channelPoints int) bool {
        if !c.CanUpgradeStat(statType, channelPoints) {
                return false
        }
        
        // Each stat upgrade costs channelPoints and increases stat by 1
        // This can be made more complex later (e.g., increasing costs)
        pointsPerUpgrade := 100 // Base cost for stat upgrade
        upgradeAmount := channelPoints / pointsPerUpgrade
        
        if upgradeAmount <= 0 {
                return false
        }
        
        switch statType {
        case "strength":
                c.Strength += upgradeAmount
        case "agility":
                c.Agility += upgradeAmount
        case "vitality":
                c.Vitality += upgradeAmount
        case "intelligence":
                c.Intelligence += upgradeAmount
        default:
                return false
        }
        
        c.ChannelPointsSpent += channelPoints
        return true
}

// GetNextLevelExperience calculates experience needed for next level
func (c *Character) GetNextLevelExperience() int {
        // Simple leveling formula: level * 1000
        return c.Level * 1000
}

// AddExperience adds experience and handles level ups
func (c *Character) AddExperience(exp int) bool {
        c.Experience += exp
        leveledUp := false
        
        // Check for level up
        for c.Experience >= c.GetNextLevelExperience() {
                c.Experience -= c.GetNextLevelExperience()
                c.Level++
                leveledUp = true
                
                // Bonus stats on level up
                c.Strength += 1
                c.Agility += 1
                c.Vitality += 1
                c.Intelligence += 1
        }
        
        return leveledUp
}