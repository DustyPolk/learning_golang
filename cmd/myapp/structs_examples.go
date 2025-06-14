package main

import (
	"fmt"
	"math/rand"
)

// This is a simple example of a struct in Go.
// A struct is a collection of fields.
// A field is a variable that is part of a struct.
// A field is a variable that is part of a struct.
type Person struct {
	Name         string
	Age          int
	City         string
	Sport        string
	Strength     int
	Speed        int
	Agility      int
	Endurance    int
	Intelligence int
	Creativity   int
	Teamwork     int
	Leadership   int
	Health       int
}

// This is a method of the Person struct.
// A method is a function that is part of a struct.
func (p Person) String() string {
	return fmt.Sprintf("%v (%v years old) from %v, %v", p.Name, p.Age, p.City, p.Sport)
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
		death(p)
	}
}

func death(p *Person) {
	p.Health = 0
	fmt.Println(p.Name, "is dead!")
}

func fightToDeath(p1 *Person, p2 *Person) {
	fmt.Printf("Fight to the death begins: %s vs %s!\n", p1.Name, p2.Name)

	for p1.Health > 0 && p2.Health > 0 {
		// Player 1 attacks Player 2
		if p1.Health > 0 {
			damage := calculateDamage(*p1)
			fmt.Printf("%s attacks %s for %d damage!\n", p1.Name, p2.Name, damage)
			takeDamage(p2, damage)
			fmt.Printf("%s health: %d\n", p2.Name, p2.Health)
		}

		// Player 2 attacks Player 1 (if still alive)
		if p2.Health > 0 {
			damage := calculateDamage(*p2)
			fmt.Printf("%s attacks %s for %d damage!\n", p2.Name, p1.Name, damage)
			takeDamage(p1, damage)
			fmt.Printf("%s health: %d\n", p1.Name, p1.Health)
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
		Name:         "Allison",
		Age:          24,
		City:         "Denver",
		Sport:        "Powerlifter",
		Strength:     80,
		Speed:        70,
		Agility:      80,
		Endurance:    80,
		Intelligence: 100,
		Creativity:   80,
		Teamwork:     80,
		Leadership:   80,
		Health:       100,
	}

	// This is a struct literal.
	// A struct literal is a literal that is used to create a struct.
	p2 := Person{
		Name:         "Dustin",
		Age:          20,
		City:         "Denver",
		Sport:        "Football",
		Strength:     100,
		Speed:        100,
		Agility:      100,
		Endurance:    100,
		Intelligence: 80,
		Creativity:   80,
		Teamwork:     80,
		Leadership:   80,
		Health:       100,
	}

	// This is a function call.
	// A function call is a function that is called.
	takeDamage(&p1, 10)
	takeDamage(&p2, 10)
	fmt.Println(p1.Health)
	fmt.Println(p2.Health)
	fightToDeath(&p1, &p2)
}
