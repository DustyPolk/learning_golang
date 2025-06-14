package combat

import (
	"fmt"
	"math/rand"

	"github.com/dustypolk/learning_golang/internal/dice"
)

// Person represents a character in combat
type Person struct {
	Name     string
	Strength int
	Health   int
	// New defensive stats
	Defense int // Reduces incoming damage
	Armor   int // Additional damage reduction
	Dodge   int // Chance to avoid damage completely (1-100)
}

// CalculateDamage determines how much damage a person deals
func CalculateDamage(p Person) int {
	attackRoll := dice.D20()
	baseDamage := dice.D6()
	strengthMod := p.Strength / 20
	totalDamage := baseDamage + strengthMod

	fmt.Printf("DEBUG: %s attacks with roll %d, base damage %d, strength mod %d, total %d\n",
		p.Name, attackRoll, baseDamage, strengthMod, totalDamage)

	if attackRoll == 1 {
		return 0 // Critical miss
	}

	if attackRoll == 20 {
		damage := dice.RollMultiple(2, 6)
		critDamage := damage * 2
		fmt.Printf("DEBUG: Critical hit! %d * 2 = %d damage\n", damage, critDamage)
		return critDamage
	}

	// Scale damage based on attack roll (1-20)
	// This makes higher rolls do more damage
	damageMultiplier := float64(attackRoll) / 20.0
	scaledDamage := int(float64(totalDamage) * damageMultiplier)

	fmt.Printf("DEBUG: Attack roll %d scales damage by %.2f: %d -> %d\n",
		attackRoll, damageMultiplier, totalDamage, scaledDamage)

	return scaledDamage
}

// TakeDamage reduces a person's health by the given amount, accounting for defense
func TakeDamage(p *Person, damage int, silent bool) {
	// Check for dodge
	dodgeRoll := rand.Intn(100)
	fmt.Printf("DEBUG: %s dodge roll: %d vs %d\n", p.Name, dodgeRoll, p.Dodge)
	if dodgeRoll < p.Dodge {
		if !silent {
			fmt.Printf("%s dodges the attack!\n", p.Name)
		}
		return
	}

	// Calculate damage reduction - increased effectiveness
	damageReduction := (p.Defense / 10) + (p.Armor / 5) // Doubled the effectiveness
	fmt.Printf("DEBUG: %s defense reduction: %d (defense %d/10 + armor %d/5)\n",
		p.Name, damageReduction, p.Defense, p.Armor)

	if damageReduction > 0 {
		damage = damage - damageReduction
		if damage < 1 {
			damage = 1 // Minimum 1 damage
		}
		if !silent {
			fmt.Printf("%s's defense reduces damage by %d!\n", p.Name, damageReduction)
		}
	}

	p.Health -= damage
	fmt.Printf("DEBUG: %s takes %d damage, health now %d\n", p.Name, damage, p.Health)

	if p.Health <= 0 {
		p.Health = 0
		if !silent {
			fmt.Println(p.Name, "is dead!")
		}
	}
}

// DecideWhoAttacksFirst determines which person attacks first
func DecideWhoAttacksFirst(p1 *Person, p2 *Person) *Person {
	p1Rolls := 0
	p2Rolls := 0

	for range 3 {
		p1Roll := dice.D20()
		p2Roll := dice.D20()

		if p1Roll > p2Roll {
			p1Rolls++
		} else if p2Roll > p1Roll {
			p2Rolls++
		}
	}

	if p1Rolls > p2Rolls {
		return p1
	}
	return p2
}

// FightToDeath handles a fight between two people
func FightToDeath(p1 *Person, p2 *Person, silent ...bool) {
	isSilent := len(silent) > 0 && silent[0]

	attacker := DecideWhoAttacksFirst(p1, p2)
	if !isSilent {
		fmt.Printf("%s wins the roll! Gets to attack first!\n", attacker.Name)
	}

	for p1.Health > 0 && p2.Health > 0 {
		// First attacker goes
		if attacker == p1 && p1.Health > 0 {
			damage := CalculateDamage(*p1)
			if !isSilent {
				fmt.Printf("%s attacks %s for %d damage!\n", p1.Name, p2.Name, damage)
			}
			TakeDamage(p2, damage, isSilent)
			if !isSilent {
				fmt.Printf("%s health: %d\n", p2.Name, p2.Health)
			}
		} else if attacker == p2 && p2.Health > 0 {
			damage := CalculateDamage(*p2)
			if !isSilent {
				fmt.Printf("%s attacks %s for %d damage!\n", p2.Name, p1.Name, damage)
			}
			TakeDamage(p1, damage, isSilent)
			if !isSilent {
				fmt.Printf("%s health: %d\n", p1.Name, p1.Health)
			}
		}

		// Second attacker goes (if still alive)
		if attacker == p1 && p2.Health > 0 {
			damage := CalculateDamage(*p2)
			if !isSilent {
				fmt.Printf("%s attacks %s for %d damage!\n", p2.Name, p1.Name, damage)
			}
			TakeDamage(p1, damage, isSilent)
			if !isSilent {
				fmt.Printf("%s health: %d\n", p1.Name, p1.Health)
			}
		} else if attacker == p2 && p1.Health > 0 {
			damage := CalculateDamage(*p1)
			if !isSilent {
				fmt.Printf("%s attacks %s for %d damage!\n", p1.Name, p2.Name, damage)
			}
			TakeDamage(p2, damage, isSilent)
			if !isSilent {
				fmt.Printf("%s health: %d\n", p2.Name, p2.Health)
			}
		}

		if !isSilent {
			fmt.Println("---")
		}
	}

	if !isSilent {
		if p1.Health > 0 {
			fmt.Printf("%s wins the fight to the death!\n", p1.Name)
		} else {
			fmt.Printf("%s wins the fight to the death!\n", p2.Name)
		}
	}
}
