package main

import (
	"math/rand"
	"time"
	"github.com/dustypolk/learning_golang/internal/game"
	"github.com/dustypolk/learning_golang/internal/graphics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize random seed for combat
	rand.Seed(time.Now().UnixNano())
	
	g := game.NewGame()
	renderer := graphics.NewRenderer(g)

	renderer.Initialize()
	defer renderer.Close()

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()
		
		g.HandleInput()
		g.Update(deltaTime)
		renderer.Draw()
	}
}