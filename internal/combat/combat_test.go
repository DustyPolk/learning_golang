package combat

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestFightWinRates(t *testing.T) {
	// Initialize random seed for true randomness
	rand.Seed(time.Now().UnixNano())
	
	// Create test characters
	allison := Person{
		Name:     "Allison",
		Strength: 80,
		Health:   110,
		Defense:  25,
		Armor:    20,
		Dodge:    20,
	}

	dustin := Person{
		Name:     "Dustin",
		Strength: 100,
		Health:   100,
		Defense:  20,
		Armor:    15,
		Dodge:    10,
	}

	// First run a single fight with detailed output
	fmt.Println("\nRunning a single detailed fight:")
	allisonCopy := allison
	dustinCopy := dustin
	FightToDeath(&allisonCopy, &dustinCopy)

	// Now run 100 fights silently
	fmt.Println("\nRunning 100 fights for analysis:")
	allisonWins := 0
	dustinWins := 0

	for i := 0; i < 100; i++ {
		// Reset health for each fight
		allisonCopy = allison
		dustinCopy = dustin

		FightToDeath(&allisonCopy, &dustinCopy, true)

		if allisonCopy.Health > 0 {
			allisonWins++
		} else {
			dustinWins++
		}
	}

	// Calculate win rates
	allisonWinRate := float64(allisonWins) / 100.0 * 100
	dustinWinRate := float64(dustinWins) / 100.0 * 100

	// Print analysis
	fmt.Printf("\nFight Analysis (100 fights):\n")
	fmt.Printf("Allison wins: %d (%.1f%%)\n", allisonWins, allisonWinRate)
	fmt.Printf("Dustin wins: %d (%.1f%%)\n", dustinWins, dustinWinRate)

	// Warn if win rates are heavily skewed
	if allisonWinRate > 70 || dustinWinRate > 70 {
		t.Logf("Warning: Win rate is heavily skewed (>70%%)")
	}
}
