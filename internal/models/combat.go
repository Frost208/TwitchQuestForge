package models

import (
        "fmt"
        "math/rand"
        "time"
)

// CombatLog represents a fight between two characters
type CombatLog struct {
        ID            int       `json:"id" db:"id"`
        AttackerID    int       `json:"attacker_id" db:"attacker_id"`
        DefenderID    int       `json:"defender_id" db:"defender_id"`
        WinnerID      int       `json:"winner_id" db:"winner_id"`
        AttackerPower int       `json:"attacker_power" db:"attacker_power"`
        DefenderPower int       `json:"defender_power" db:"defender_power"`
        CombatLogText string    `json:"combat_log" db:"combat_log"`
        CreatedAt     time.Time `json:"created_at" db:"created_at"`
        
        // Populated fields
        Attacker *Character `json:"attacker,omitempty"`
        Defender *Character `json:"defender,omitempty"`
        Winner   *Character `json:"winner,omitempty"`
}

// CombatRequest represents a request to start combat
type CombatRequest struct {
        DefenderUsername string `json:"defender_username" binding:"required"`
}

// CombatResult represents the result of a combat encounter
type CombatResult struct {
        Winner           *Character `json:"winner"`
        Loser            *Character `json:"loser"`
        AttackerPower    int        `json:"attacker_power"`
        DefenderPower    int        `json:"defender_power"`
        CombatLog        string     `json:"combat_log"`
        ExperienceGained int        `json:"experience_gained"`
        RewardItems      []Item     `json:"reward_items,omitempty"`
}

// SimulateCombat simulates a combat encounter between two characters
func SimulateCombat(attacker, defender *Character) *CombatResult {
        // Calculate combat powers
        attackerPower := attacker.CalculateCombatPower()
        defenderPower := defender.CalculateCombatPower()
        
        // Add some randomness (±20% variation)
        attackerRoll := attackerPower + (attackerPower * (rand.Intn(40) - 20) / 100)
        defenderRoll := defenderPower + (defenderPower * (rand.Intn(40) - 20) / 100)
        
        var winner, loser *Character
        var combatLog string
        
        if attackerRoll > defenderRoll {
                winner = attacker
                loser = defender
                combatLog = fmt.Sprintf("%s besiegt %s! (%d vs %d Kampfkraft)", 
                        attacker.Username, defender.Username, attackerRoll, defenderRoll)
        } else {
                winner = defender
                loser = attacker
                combatLog = fmt.Sprintf("%s verteidigt sich erfolgreich gegen %s! (%d vs %d Kampfkraft)", 
                        defender.Username, attacker.Username, defenderRoll, attackerRoll)
        }
        
        // Calculate experience reward (based on opponent's level)
        experienceGained := loser.Level * 50
        if winner == defender {
                experienceGained = int(float64(experienceGained) * 1.2) // Defender bonus
        }
        
        return &CombatResult{
                Winner:           winner,
                Loser:            loser,
                AttackerPower:    attackerPower,
                DefenderPower:    defenderPower,
                CombatLog:        combatLog,
                ExperienceGained: experienceGained,
        }
}