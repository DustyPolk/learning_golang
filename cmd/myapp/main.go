package main

import (
	"fmt"

	"github.com/dustypolk/learning_golang/internal/combat"
)

func main() {
	p1 := combat.Person{
		Name:     "Allison",
		Strength: 90,
		Health:   150,
		Defense:  80,
		Armor:    80,
		Dodge:    40,
	}

	p2 := combat.Person{
		Name:     "Dustin",
		Strength: 100,
		Health:   100,
		Defense:  20,
		Armor:    20,
		Dodge:    5,
	}

	// Run a single fight with output
	fmt.Println("Running a single fight:")
	combat.FightToDeath(&p1, &p2)

	// Run 100 fights silently and analyze results
	fmt.Println("\nRunning 100 fights for analysis:")
	allisonWins := 0
	dustinWins := 0

	for i := 0; i < 100; i++ {
		// Reset health for each fight
		p1.Health = 150
		p2.Health = 100

		combat.FightToDeath(&p1, &p2, true)

		if p1.Health > 0 {
			allisonWins++
		} else {
			dustinWins++
		}
	}

	// Print analysis
	fmt.Printf("\nFight Analysis (100 fights):\n")
	fmt.Printf("Allison wins: %d (%.1f%%)\n", allisonWins, float64(allisonWins)/100*100)
	fmt.Printf("Dustin wins: %d (%.1f%%)\n", dustinWins, float64(dustinWins)/100*100)
}
