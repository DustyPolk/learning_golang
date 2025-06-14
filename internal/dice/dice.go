package dice

import "math/rand"

// Roll rolls a die with the specified number of sides
func Roll(sides int) int {
	return rand.Intn(sides) + 1
}

// RollMultiple rolls multiple dice and returns the sum
func RollMultiple(count, sides int) int {
	total := 0
	for range count {
		total += Roll(sides)
	}
	return total
}

// D20 rolls a 20-sided die
func D20() int {
	return Roll(20)
}

// D6 rolls a 6-sided die
func D6() int {
	return Roll(6)
}
