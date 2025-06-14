package combat

import (
	"fmt"
	"testing"
)

func TestFightWinRates(t *testing.T) {
	const numFights = 100
	allisonWins := 0
	dustinWins := 0

	for i := 0; i < numFights; i++ {
		// Create fresh fighters for each battle
		p1 := Person{
			Name:     "Allison",
			Strength: 80,
			Health:   100,
			Defense:  50,
			Armor:    30,
			Dodge:    15,
		}

		p2 := Person{
			Name:     "Dustin",
			Strength: 100,
			Health:   120,
			Defense:  40,
			Armor:    60,
			Dodge:    10,
		}

		// Fight to the death
		FightToDeath(&p1, &p2)

		// Record winner
		if p1.Health > 0 {
			allisonWins++
		} else {
			dustinWins++
		}
	}

	// Calculate win rates
	allisonWinRate := float64(allisonWins) / float64(numFights) * 100
	dustinWinRate := float64(dustinWins) / float64(numFights) * 100

	// Print results
	fmt.Printf("\nFight Analysis (100 fights):\n")
	fmt.Printf("Allison wins: %d (%.1f%%)\n", allisonWins, allisonWinRate)
	fmt.Printf("Dustin wins: %d (%.1f%%)\n", dustinWins, dustinWinRate)

	// Optional: Add test assertions
	if allisonWinRate > 70 || dustinWinRate > 70 {
		t.Logf("Warning: Win rate is heavily skewed (>70%%)")
	}
}
