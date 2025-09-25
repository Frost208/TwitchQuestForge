package services

import (
        "fmt"
        "math/rand"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/storage"
)

// CombatService handles combat-related operations
type CombatService struct{}

// NewCombatService creates a new combat service
func NewCombatService() *CombatService {
        return &CombatService{}
}

// StartCombat initiates combat between two characters
func (cs *CombatService) StartCombat(attackerID, defenderID int) (*models.CombatResult, error) {
        if database.DB == nil {
                // Use memory storage for testing
                return cs.startCombatMemory(attackerID, defenderID)
        }

        // Get both characters
        charService := NewCharacterService()
        attacker, err := charService.GetCharacterByID(attackerID)
        if err != nil {
                return nil, fmt.Errorf("failed to get attacker: %v", err)
        }
        if attacker == nil {
                return nil, fmt.Errorf("attacker not found")
        }

        defender, err := charService.GetCharacterByID(defenderID)
        if err != nil {
                return nil, fmt.Errorf("failed to get defender: %v", err)
        }
        if defender == nil {
                return nil, fmt.Errorf("defender not found")
        }

        // Calculate combat power for both characters
        attackerPower := attacker.CalculateCombatPower()
        defenderPower := defender.CalculateCombatPower()

        // Simulate combat with some randomness
        attackerChance := float64(attackerPower) / float64(attackerPower+defenderPower)
        randomRoll := rand.Float64()

        winner := attacker
        loser := defender
        if randomRoll > attackerChance {
                winner = defender
                loser = attacker
        }

        // Calculate experience and rewards
        experienceGained := 50 + (loser.Level * 10)
        channelPointsReward := 25 + (loser.Level * 5)

        // Update winner's stats
        winner.Experience += experienceGained
        winner.ChannelPointsSpent -= channelPointsReward // Give channel points as reward

        // Check for level up
        expForNextLevel := winner.Level * 100
        if winner.Experience >= expForNextLevel {
                winner.Level++
                winner.Experience -= expForNextLevel
        }

        // Save changes
        err = charService.UpdateCharacter(winner)
        if err != nil {
                return nil, fmt.Errorf("failed to update winner: %v", err)
        }

        // Create combat log
        combatResult := &models.CombatResult{
                Winner:           winner,
                Loser:            loser,
                AttackerPower:    attackerPower,
                DefenderPower:    defenderPower,
                CombatLog:        fmt.Sprintf("%s defeated %s! Experience gained: %d", winner.Username, loser.Username, experienceGained),
                ExperienceGained: experienceGained,
                RewardItems:      []models.Item{}, // No item rewards for now
        }

        err = cs.logCombat(combatResult)
        if err != nil {
                return nil, fmt.Errorf("failed to log combat: %v", err)
        }

        return combatResult, nil
}

// GetCombatHistory retrieves recent combat history
func (cs *CombatService) GetCombatHistory(limit int) ([]models.CombatLog, error) {
        if database.DB == nil {
                return storage.Memory.GetCombatHistory(limit)
        }

        query := `
                SELECT cl.id, cl.attacker_id, cl.defender_id, cl.winner_id, 
                       cl.attacker_power, cl.defender_power, cl.combat_log, cl.created_at
                FROM combat_logs cl
                ORDER BY cl.created_at DESC
                LIMIT ?`

        rows, err := database.DB.Query(query, limit)
        if err != nil {
                return nil, fmt.Errorf("failed to get combat history: %v", err)
        }
        defer rows.Close()

        var combatLogs []models.CombatLog
        for rows.Next() {
                var log models.CombatLog
                err := rows.Scan(
                        &log.ID, &log.AttackerID, &log.DefenderID, &log.WinnerID,
                        &log.AttackerPower, &log.DefenderPower, &log.CombatLogText,
                        &log.CreatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan combat log: %v", err)
                }
                combatLogs = append(combatLogs, log)
        }

        return combatLogs, nil
}

// logCombat records a combat result to the database
func (cs *CombatService) logCombat(result *models.CombatResult) error {
        query := `
                INSERT INTO combat_logs (attacker_id, defender_id, winner_id, attacker_power, 
                                        defender_power, combat_log)
                VALUES (?, ?, ?, ?, ?, ?)`

        // Determine attacker/defender from winner/loser
        attackerID := result.Winner.ID  // Simplified logic
        defenderID := result.Loser.ID

        _, err := database.DB.Exec(query,
                attackerID,
                defenderID,
                result.Winner.ID,
                result.AttackerPower,
                result.DefenderPower,
                result.CombatLog,
        )

        return err
}