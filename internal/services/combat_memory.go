package services

import (
	"fmt"
	"math/rand"
	"twitch-rpg/internal/models"
	"twitch-rpg/internal/storage"
)

// startCombatMemory handles combat using memory storage
func (cs *CombatService) startCombatMemory(attackerID, defenderID int) (*models.CombatResult, error) {
	// Get both characters from memory
	attacker, err := storage.Memory.GetCharacterByID(attackerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attacker: %v", err)
	}
	if attacker == nil {
		return nil, fmt.Errorf("attacker not found")
	}

	defender, err := storage.Memory.GetCharacterByID(defenderID)
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

	// Update winner's stats
	winner.Experience += experienceGained

	// Check for level up
	expForNextLevel := winner.Level * 100
	if winner.Experience >= expForNextLevel {
		winner.Level++
		winner.Experience -= expForNextLevel
	}

	// Save changes to memory
	err = storage.Memory.UpdateCharacter(winner)
	if err != nil {
		return nil, fmt.Errorf("failed to update winner: %v", err)
	}

	// Create combat result
	combatResult := &models.CombatResult{
		Winner:           winner,
		Loser:            loser,
		AttackerPower:    attackerPower,
		DefenderPower:    defenderPower,
		CombatLog:        fmt.Sprintf("%s defeated %s! Experience gained: %d", winner.Username, loser.Username, experienceGained),
		ExperienceGained: experienceGained,
		RewardItems:      []models.Item{}, // No item rewards for now
	}

	// Log combat to memory
	combatLog := models.CombatLog{
		AttackerID:    attackerID,
		DefenderID:    defenderID,
		WinnerID:      winner.ID,
		AttackerPower: attackerPower,
		DefenderPower: defenderPower,
		CombatLogText: combatResult.CombatLog,
	}

	err = storage.Memory.AddCombatLog(combatLog)
	if err != nil {
		return nil, fmt.Errorf("failed to log combat: %v", err)
	}

	return combatResult, nil
}