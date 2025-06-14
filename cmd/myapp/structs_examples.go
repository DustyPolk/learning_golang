package main

import (
	"fmt"
	"math/rand"
)

// This is a simple example of a struct in Go.
// A struct is a collection of fields.
// A field is a variable that is part of a struct.
type Person struct {
	Name     string
	Strength int
	Health   int
}

// This is a function that is not part of the Person struct.
// A function is a collection of statements that perform a task.
// This function calculates the damage a player deals to another player.
func calculateDamage(p Person) int {
	// Roll a d20 for attack roll
	attackRoll := rand.Intn(20) + 1

	// If attack roll is 1, it's a critical miss
	if attackRoll == 1 {
		return 0
	}

	// If attack roll is 20, it's a critical hit
	if attackRoll == 20 {
		// Roll 2d6 for critical hit damage
		damage := (rand.Intn(6) + 1) + (rand.Intn(6) + 1)
		return damage * 2
	}

	// Normal hit - roll 1d6 for damage
	damage := rand.Intn(6) + 1
	// Add strength modifier (strength/10 rounded down)
	strengthMod := p.Strength / 10
	return damage + strengthMod
}

func takeDamage(p *Person, damage int) {
	p.Health -= damage
	if p.Health <= 0 {
		p.Health = 0
		fmt.Println(p.Name, "is dead!")
	}
}

func decideWhoAttacksFirst(p1 *Person, p2 *Person) *Person {
	// Roll 3 times for each player
	p1Rolls := 0
	p2Rolls := 0

	for range 3 {
		p1Roll := rand.Intn(20) + 1
		p2Roll := rand.Intn(20) + 1

		if p1Roll > p2Roll {
			p1Rolls++
		} else if p2Roll > p1Roll {
			p2Rolls++
		}
	}

	// Return the player who won more rolls
	if p1Rolls > p2Rolls {
		return p1
	}
	return p2
}

func fightToDeath(p1 *Person, p2 *Person) {
	attacker := decideWhoAttacksFirst(p1, p2)
	fmt.Printf("%s wins the roll! Gets to attack first!\n", attacker.Name)

	for p1.Health > 0 && p2.Health > 0 {
		// First attacker goes
		if attacker == p1 && p1.Health > 0 {
			damage := calculateDamage(*p1)
			fmt.Printf("%s attacks %s for %d damage!\n", p1.Name, p2.Name, damage)
			takeDamage(p2, damage)
			fmt.Printf("%s health: %d\n", p2.Name, p2.Health)
		} else if attacker == p2 && p2.Health > 0 {
			damage := calculateDamage(*p2)
			fmt.Printf("%s attacks %s for %d damage!\n", p2.Name, p1.Name, damage)
			takeDamage(p1, damage)
			fmt.Printf("%s health: %d\n", p1.Name, p1.Health)
		}

		// Second attacker goes (if still alive)
		if attacker == p1 && p2.Health > 0 {
			damage := calculateDamage(*p2)
			fmt.Printf("%s attacks %s for %d damage!\n", p2.Name, p1.Name, damage)
			takeDamage(p1, damage)
			fmt.Printf("%s health: %d\n", p1.Name, p1.Health)
		} else if attacker == p2 && p1.Health > 0 {
			damage := calculateDamage(*p1)
			fmt.Printf("%s attacks %s for %d damage!\n", p1.Name, p2.Name, damage)
			takeDamage(p2, damage)
			fmt.Printf("%s health: %d\n", p2.Name, p2.Health)
		}

		fmt.Println("---")
	}

	if p1.Health > 0 {
		fmt.Printf("%s wins the fight to the death!\n", p1.Name)
	} else {
		fmt.Printf("%s wins the fight to the death!\n", p2.Name)
	}
}

// This is the main function.
// The main function is the entry point of the program.
func main() {
	// This is a struct literal.
	// A struct literal is a literal that is used to create a struct.
	p1 := Person{
		Name:     "Allison",
		Strength: 80,
		Health:   100,
	}

	// This is a struct literal.
	// A struct literal is a literal that is used to create a struct.
	p2 := Person{
		Name:     "Dustin",
		Strength: 100,
		Health:   100,
	}

	// This is a function call.
	// A function call is a function that is called.
	fightToDeath(&p1, &p2)
}
