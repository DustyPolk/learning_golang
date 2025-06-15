package graphics

import (
	"fmt"
	"github.com/dustypolk/learning_golang/internal/combat"
	"github.com/dustypolk/learning_golang/internal/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

type Renderer struct {
	game *game.Game
}

func NewRenderer(g *game.Game) *Renderer {
	return &Renderer{game: g}
}

func (r *Renderer) Initialize() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Go Combat Arena")
	rl.SetTargetFPS(60)
}

func (r *Renderer) Close() {
	rl.CloseWindow()
}

func (r *Renderer) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)

	switch r.game.State {
	case game.StateMenu:
		r.drawMenu()
	case game.StateBattle, game.StateActionMenu, game.StateTargetSelect, game.StateSpellSelect, game.StateSpellTarget:
		r.drawBattle()
	case game.StateVictory:
		r.drawVictoryScreen()
	case game.StateDefeat:
		r.drawDefeatScreen()
	}

	rl.EndDrawing()
}

func (r *Renderer) drawMenu() {
	centerX := int32(ScreenWidth / 2)
	centerY := int32(ScreenHeight / 2)

	title := "GO COMBAT ARENA"
	titleSize := int32(60)
	titleWidth := rl.MeasureText(title, titleSize)
	rl.DrawText(title, centerX-titleWidth/2, centerY-100, titleSize, rl.RayWhite)

	instructions := "Press SPACE to Start Battle"
	instrSize := int32(30)
	instrWidth := rl.MeasureText(instructions, instrSize)
	rl.DrawText(instructions, centerX-instrWidth/2, centerY, instrSize, rl.LightGray)
}

func (r *Renderer) drawBattle() {
	r.drawBackground()
	r.drawCharacters()
	r.drawBattleLog()
	r.drawTurnOrder()
	r.drawAnimations()
	
	// Draw different UI based on state
	switch r.game.State {
	case game.StateActionMenu:
		r.drawActionMenu()
	case game.StateTargetSelect:
		r.drawTargetSelection()
	case game.StateSpellSelect:
		r.drawSpellSelection()
	case game.StateSpellTarget:
		r.drawSpellTargetSelection()
	default:
		r.drawTurnInfo()
	}
}

func (r *Renderer) drawBackground() {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight/2, rl.SkyBlue)
	rl.DrawRectangle(0, ScreenHeight/2, ScreenWidth, ScreenHeight/2, rl.Brown)
}

func (r *Renderer) drawCharacters() {
	// Draw player party (left side)
	for i, char := range r.game.PlayerParty {
		if char.Health > 0 {
			baseX := int32(100)
			baseY := int32(200 + i*120)
			
			// Check for melee animation offset
			offsetX := r.getMeleeOffset(char)
			
			color := rl.Blue
			if char.Name == "Mage" {
				color = rl.Purple
			} else if char.Name == "Rogue" {
				color = rl.Green
			}
			
			r.drawCharacter(char, baseX+offsetX, baseY, color)
		}
	}
	
	// Draw enemy party (right side)
	for i, char := range r.game.EnemyParty {
		if char.Health > 0 {
			baseX := int32(ScreenWidth - 200)
			baseY := int32(250 + i*120)
			
			// Check for melee animation offset
			offsetX := r.getMeleeOffset(char)
			
			color := rl.Red
			if char.Name == "Goblin" {
				color = rl.Orange
			}
			
			r.drawCharacter(char, baseX+offsetX, baseY, color)
		}
	}
}

func (r *Renderer) getMeleeOffset(character *combat.Person) int32 {
	for _, action := range r.game.ActionQueue {
		if action.Type == game.ActionAttack && action.AttackType == game.AttackMelee && action.Attacker == character {
			elapsedTime := r.game.AnimationTimer - action.StartTime
			animDuration := float32(0.8)
			progress := elapsedTime / animDuration
			if progress > 1 {
				progress = 1
			}
			
			// Calculate movement toward enemy side
			targetOffset := float32(300)
			isPlayerChar := r.isPlayerCharacter(character)
			
			if progress < 0.5 {
				// Moving to target
				t := progress * 2
				if isPlayerChar {
					return int32(targetOffset * t) // Move right
				} else {
					return -int32(targetOffset * t) // Move left
				}
			} else {
				// Moving back
				t := (progress - 0.5) * 2
				if isPlayerChar {
					return int32(targetOffset * (1 - t))
				} else {
					return -int32(targetOffset * (1 - t))
				}
			}
		}
	}
	return 0
}

func (r *Renderer) isPlayerCharacter(character *combat.Person) bool {
	for _, char := range r.game.PlayerParty {
		if char == character {
			return true
		}
	}
	return false
}

func (r *Renderer) drawCharacter(character *combat.Person, x, y int32, color rl.Color) {
	charWidth := int32(100)
	charHeight := int32(150)

	rl.DrawRectangle(x, y, charWidth, charHeight, color)
	
	nameSize := int32(20)
	nameWidth := rl.MeasureText(character.Name, nameSize)
	rl.DrawText(character.Name, x+charWidth/2-nameWidth/2, y-30, nameSize, rl.White)

	r.drawHealthBar(character, x, y-50, charWidth)
	r.drawMPBar(character, x, y-25, charWidth)

	if r.game.CurrentCharacter == character {
		// Draw selection indicator
		rl.DrawTriangle(
			rl.NewVector2(float32(x+charWidth/2), float32(y-70)),
			rl.NewVector2(float32(x+charWidth/2-10), float32(y-85)),
			rl.NewVector2(float32(x+charWidth/2+10), float32(y-85)),
			rl.Yellow,
		)
		
		// Draw glow effect around current character
		rl.DrawRectangleLines(x-5, y-5, charWidth+10, charHeight+10, rl.Yellow)
	}
}

func (r *Renderer) drawHealthBar(character *combat.Person, x, y, width int32) {
	barHeight := int32(20)
	maxHealth := 110
	if character.Name == "Player" {
		maxHealth = 100
	}
	
	healthPercent := float32(character.Health) / float32(maxHealth)
	if healthPercent < 0 {
		healthPercent = 0
	}

	rl.DrawRectangle(x, y, width, barHeight, rl.DarkGray)
	rl.DrawRectangle(x, y, int32(float32(width)*healthPercent), barHeight, rl.Green)
	rl.DrawRectangleLines(x, y, width, barHeight, rl.Black)

	healthText := fmt.Sprintf("%d/%d", character.Health, maxHealth)
	textSize := int32(16)
	textWidth := rl.MeasureText(healthText, textSize)
	rl.DrawText(healthText, x+width/2-textWidth/2, y+2, textSize, rl.White)
}

func (r *Renderer) drawMPBar(character *combat.Person, x, y, width int32) {
	barHeight := int32(15)
	
	mpPercent := float32(character.MP) / float32(character.MaxMP)
	if mpPercent < 0 {
		mpPercent = 0
	}

	rl.DrawRectangle(x, y, width, barHeight, rl.DarkGray)
	rl.DrawRectangle(x, y, int32(float32(width)*mpPercent), barHeight, rl.Blue)
	rl.DrawRectangleLines(x, y, width, barHeight, rl.Black)

	mpText := fmt.Sprintf("MP: %d/%d", character.MP, character.MaxMP)
	textSize := int32(12)
	textWidth := rl.MeasureText(mpText, textSize)
	rl.DrawText(mpText, x+width/2-textWidth/2, y+1, textSize, rl.White)
}

func (r *Renderer) drawTurnOrder() {
	if len(r.game.TurnOrder) == 0 {
		return
	}
	
	// Draw turn order panel
	panelX := int32(ScreenWidth - 300)
	panelY := int32(50)
	panelWidth := int32(250)
	panelHeight := int32(150)
	
	rl.DrawRectangle(panelX, panelY, panelWidth, panelHeight, rl.NewColor(0, 0, 0, 200))
	rl.DrawRectangleLines(panelX, panelY, panelWidth, panelHeight, rl.White)
	
	// Title
	rl.DrawText("Turn Order", panelX+10, panelY+5, 18, rl.White)
	
	// Show next 5 turns
	for i := 0; i < 5 && i < len(r.game.TurnOrder); i++ {
		turnIndex := (r.game.CurrentTurnIndex + i) % len(r.game.TurnOrder)
		character := r.game.TurnOrder[turnIndex]
		
		if character.Health <= 0 {
			continue
		}
		
		y := panelY + 25 + int32(i*20)
		
		// Highlight current turn
		if i == 0 {
			rl.DrawRectangle(panelX+5, y-2, panelWidth-10, 18, rl.NewColor(100, 100, 0, 100))
		}
		
		// Character name with color coding
		textColor := rl.LightGray
		if r.isPlayerCharacter(character) {
			textColor = rl.Green
		} else {
			textColor = rl.Red
		}
		
		turnText := fmt.Sprintf("%d. %s", i+1, character.Name)
		if i == 0 {
			turnText += " ←"
			textColor = rl.Yellow
		}
		
		rl.DrawText(turnText, panelX+10, y, 16, textColor)
	}
}

func (r *Renderer) drawActionMenu() {
	uiY := int32(ScreenHeight - 250)
	uiWidth := int32(400)
	uiHeight := int32(200)
	uiX := int32(ScreenWidth - uiWidth - 50)
	
	// Draw menu background
	rl.DrawRectangle(uiX, uiY, uiWidth, uiHeight, rl.NewColor(20, 20, 50, 240))
	rl.DrawRectangleLines(uiX, uiY, uiWidth, uiHeight, rl.White)
	
	// Draw title
	titleText := fmt.Sprintf("%s's Turn", r.game.CurrentCharacter.Name)
	titleWidth := rl.MeasureText(titleText, 24)
	rl.DrawText(titleText, uiX+uiWidth/2-titleWidth/2, uiY+10, 24, rl.White)
	
	// Draw menu options
	menuItems := []string{"Attack", "Magic", "Item", "Run"}
	for i, item := range menuItems {
		optionY := uiY + 50 + int32(i*30)
		textColor := rl.LightGray
		
		if game.BattleMenuOption(i) == r.game.SelectedMenuOption {
			// Highlight selected option
			rl.DrawRectangle(uiX+10, optionY-5, uiWidth-20, 25, rl.DarkBlue)
			textColor = rl.Yellow
		}
		
		rl.DrawText(item, uiX+20, optionY, 20, textColor)
	}
	
	// Draw instructions
	rl.DrawText("↑↓ Navigate  ENTER Select", uiX+10, uiY+uiHeight-25, 16, rl.LightGray)
}

func (r *Renderer) drawTargetSelection() {
	uiY := int32(ScreenHeight - 250)
	uiWidth := int32(400)
	uiHeight := int32(200)
	uiX := int32(ScreenWidth - uiWidth - 50)
	
	// Draw menu background
	rl.DrawRectangle(uiX, uiY, uiWidth, uiHeight, rl.NewColor(50, 20, 20, 240))
	rl.DrawRectangleLines(uiX, uiY, uiWidth, uiHeight, rl.White)
	
	// Draw title
	rl.DrawText("Select Target", uiX+20, uiY+10, 24, rl.White)
	
	// Draw enemy targets
	for i, enemy := range r.game.EnemyParty {
		if enemy.Health > 0 {
			optionY := uiY + 50 + int32(i*30)
			textColor := rl.LightGray
			
			if i == r.game.SelectedTarget {
				// Highlight selected target
				rl.DrawRectangle(uiX+10, optionY-5, uiWidth-20, 25, rl.Maroon)
				textColor = rl.Yellow
				
				// Also highlight the enemy character
				r.highlightTarget(enemy)
			}
			
			healthText := fmt.Sprintf("%s (HP: %d)", enemy.Name, enemy.Health)
			rl.DrawText(healthText, uiX+20, optionY, 20, textColor)
		}
	}
	
	// Draw instructions
	rl.DrawText("↑↓ Navigate  ENTER Attack  ESC Back", uiX+10, uiY+uiHeight-25, 16, rl.LightGray)
}

func (r *Renderer) drawTurnInfo() {
	if r.game.CurrentCharacter != nil {
		turnText := fmt.Sprintf("%s's Turn", r.game.CurrentCharacter.Name)
		textColor := rl.White
		
		if !r.isPlayerCharacter(r.game.CurrentCharacter) {
			turnText += " (Enemy)"
			textColor = rl.Red
		} else {
			textColor = rl.Green
		}
		
		// Draw turn info background
		textWidth := rl.MeasureText(turnText, 24)
		rl.DrawRectangle(40, ScreenHeight-80, textWidth+20, 35, rl.NewColor(0, 0, 0, 180))
		rl.DrawText(turnText, 50, ScreenHeight-70, 24, textColor)
		
		// Show action status
		var statusText string
		var statusColor rl.Color = rl.LightGray
		
		switch r.game.State {
		case game.StateActionMenu:
			statusText = "Choose Action..."
			statusColor = rl.Yellow
		case game.StateTargetSelect:
			statusText = "Select Target..."
			statusColor = rl.Orange
		case game.StateBattle:
			if len(r.game.ActionQueue) > 0 {
				statusText = "Executing Action..."
				statusColor = rl.SkyBlue
			} else if !r.isPlayerCharacter(r.game.CurrentCharacter) {
				statusText = "Enemy Thinking..."
				statusColor = rl.Red
			} else {
				statusText = "Ready"
				statusColor = rl.Green
			}
		}
		
		if statusText != "" {
			statusWidth := rl.MeasureText(statusText, 16)
			rl.DrawRectangle(40, ScreenHeight-40, statusWidth+10, 20, rl.NewColor(0, 0, 0, 180))
			rl.DrawText(statusText, 45, ScreenHeight-35, 16, statusColor)
		}
	}
}

func (r *Renderer) highlightTarget(target *combat.Person) {
	// Find target position and draw highlight
	for i, enemy := range r.game.EnemyParty {
		if enemy == target {
			baseX := int32(ScreenWidth - 200)
			baseY := int32(250 + i*120)
			charWidth := int32(100)
			charHeight := int32(150)
			
			// Draw pulsing highlight
			alpha := 0.5 + 0.3*float32(rl.GetTime()*4) // Pulsing effect
			if alpha > 1 {
				alpha = 1
			}
			rl.DrawRectangleLines(baseX-10, baseY-10, charWidth+20, charHeight+20, rl.Fade(rl.Red, alpha))
			break
		}
	}
}

func (r *Renderer) drawStats(character *combat.Person, x, y int32) {
	textSize := int32(16)
	lineHeight := int32(20)

	rl.DrawText(character.Name+" Stats:", x, y, textSize, rl.White)
	rl.DrawText(fmt.Sprintf("STR: %d", character.Strength), x, y+lineHeight, textSize, rl.LightGray)
	rl.DrawText(fmt.Sprintf("DEF: %d", character.Defense), x, y+lineHeight*2, textSize, rl.LightGray)
	rl.DrawText(fmt.Sprintf("DODGE: %d%%", character.Dodge), x, y+lineHeight*3, textSize, rl.LightGray)
}

func (r *Renderer) drawBattleLog() {
	logX := int32(ScreenWidth/2 - 200)
	logY := int32(50)
	logWidth := int32(400)
	logHeight := int32(200)

	rl.DrawRectangle(logX, logY, logWidth, logHeight, rl.NewColor(0, 0, 0, 180))
	rl.DrawRectangleLines(logX, logY, logWidth, logHeight, rl.White)

	textSize := int32(16)
	lineHeight := int32(20)
	
	for i, log := range r.game.BattleLog {
		rl.DrawText(log, logX+10, logY+10+int32(i)*lineHeight, textSize, rl.White)
	}
}

func (r *Renderer) drawAnimations() {
	for _, action := range r.game.ActionQueue {
		if action.Type == game.ActionAttack {
			elapsedTime := r.game.AnimationTimer - action.StartTime
			
			switch action.AttackType {
			case game.AttackMelee:
				r.drawMeleeAnimation(action, elapsedTime)
			case game.AttackRanged:
				r.drawRangedAnimation(action, elapsedTime)
			case game.AttackSpecial:
				r.drawSpecialAnimation(action, elapsedTime)
			}
		}
	}
}

func (r *Renderer) drawMeleeAnimation(action game.Action, elapsedTime float32) {
	// Character movement is now handled in drawCharacters()
	// Here we just draw the impact effects
	animDuration := float32(0.8)
	progress := elapsedTime / animDuration
	if progress > 1 {
		progress = 1
	}
	
	// Draw slash effect at impact
	if progress > 0.4 && progress < 0.6 && !action.IsDodged {
		var slashX, slashY int32
		if r.isPlayerCharacter(action.Attacker) {
			// Slash appears near enemy
			slashX = int32(ScreenWidth - 300)
			slashY = int32(ScreenHeight/2 - 25)
		} else {
			// Slash appears near player
			slashX = 300
			slashY = int32(ScreenHeight/2 - 25)
		}
		
		slashColor := rl.White
		if action.IsCritical {
			slashColor = rl.Gold
		}
		
		// Draw X-shaped slash
		slashSize := int32(40)
		rl.DrawLineEx(
			rl.NewVector2(float32(slashX-slashSize), float32(slashY-slashSize)),
			rl.NewVector2(float32(slashX+slashSize), float32(slashY+slashSize)),
			4, slashColor)
		rl.DrawLineEx(
			rl.NewVector2(float32(slashX-slashSize), float32(slashY+slashSize)),
			rl.NewVector2(float32(slashX+slashSize), float32(slashY-slashSize)),
			4, slashColor)
			
		// Add impact particles
		for i := 0; i < 5; i++ {
			offsetX := float32(i-2) * 10
			offsetY := float32(i-2) * 5
			rl.DrawCircle(slashX+int32(offsetX), slashY+int32(offsetY), 3, rl.Fade(slashColor, 0.5))
		}
	}
}

func (r *Renderer) drawRangedAnimation(action game.Action, elapsedTime float32) {
	// Projectile flies across
	progress := elapsedTime * 2
	if progress > 1 {
		progress = 1
	}

	var startX, startY, endX, endY float32
	if r.isPlayerCharacter(action.Attacker) {
		startX = 300
		startY = float32(ScreenHeight/2 - 25)
		endX = float32(ScreenWidth - 300)
		endY = startY
	} else {
		startX = float32(ScreenWidth - 300)
		startY = float32(ScreenHeight/2 - 25)
		endX = 300
		endY = startY
	}

	currentX := startX + (endX-startX)*progress
	currentY := startY + (endY-startY)*progress

	if !action.IsDodged {
		color := rl.Yellow
		if action.IsCritical {
			color = rl.Orange
		}
		rl.DrawCircle(int32(currentX), int32(currentY), 8, color)
		// Trail effect
		rl.DrawCircle(int32(currentX-10), int32(currentY), 5, rl.Fade(color, 0.5))
		rl.DrawCircle(int32(currentX-20), int32(currentY), 3, rl.Fade(color, 0.3))
	}
}

func (r *Renderer) drawSpecialAnimation(action game.Action, elapsedTime float32) {
	// Big explosion effect
	animDuration := float32(1.0)
	progress := elapsedTime / animDuration
	if progress > 1 {
		progress = 1
	}
	
	var centerX, centerY float32
	if r.isPlayerCharacter(action.Attacker) {
		centerX = float32(ScreenWidth - 300)
		centerY = float32(ScreenHeight/2 - 25)
	} else {
		centerX = 300
		centerY = float32(ScreenHeight/2 - 25)
	}
	
	if !action.IsDodged {
		// Multiple expanding circles
		for i := 0; i < 3; i++ {
			radius := progress * float32(50+i*20)
			alpha := 1.0 - progress
			color := rl.Purple
			if action.IsCritical {
				color = rl.Gold
			}
			rl.DrawCircleLines(int32(centerX), int32(centerY), radius, rl.Fade(color, alpha))
		}
		
		// Central blast
		if progress < 0.5 {
			blastRadius := progress * 60
			rl.DrawCircle(int32(centerX), int32(centerY), blastRadius, rl.Fade(rl.Purple, 0.5))
		}
	}
}

func (r *Renderer) drawSpellSelection() {
	uiY := int32(ScreenHeight - 300)
	uiWidth := int32(450)
	uiHeight := int32(250)
	uiX := int32(ScreenWidth - uiWidth - 50)
	
	// Draw menu background
	rl.DrawRectangle(uiX, uiY, uiWidth, uiHeight, rl.NewColor(20, 20, 80, 240))
	rl.DrawRectangleLines(uiX, uiY, uiWidth, uiHeight, rl.White)
	
	// Draw title
	titleText := fmt.Sprintf("%s's Spells", r.game.CurrentCharacter.Name)
	titleWidth := rl.MeasureText(titleText, 24)
	rl.DrawText(titleText, uiX+uiWidth/2-titleWidth/2, uiY+10, 24, rl.White)
	
	// Show current MP
	mpText := fmt.Sprintf("MP: %d/%d", r.game.CurrentCharacter.MP, r.game.CurrentCharacter.MaxMP)
	rl.DrawText(mpText, uiX+20, uiY+40, 18, rl.SkyBlue)
	
	// Draw spell options
	for i, spell := range r.game.CurrentSpells {
		optionY := uiY + 70 + int32(i*35)
		textColor := rl.LightGray
		
		if i == r.game.SelectedSpell {
			// Highlight selected spell
			rl.DrawRectangle(uiX+10, optionY-5, uiWidth-20, 30, rl.DarkBlue)
			textColor = rl.Yellow
		}
		
		// Check if spell is affordable
		if r.game.CurrentCharacter.MP < spell.MPCost {
			textColor = rl.DarkGray
		}
		
		spellText := fmt.Sprintf("%s (MP: %d)", spell.Name, spell.MPCost)
		rl.DrawText(spellText, uiX+20, optionY, 20, textColor)
		
		// Draw description
		rl.DrawText(spell.Description, uiX+25, optionY+18, 14, rl.LightGray)
	}
	
	// Draw instructions
	rl.DrawText("↑↓ Navigate  ENTER Select  ESC Back", uiX+10, uiY+uiHeight-25, 16, rl.LightGray)
}

func (r *Renderer) drawSpellTargetSelection() {
	uiY := int32(ScreenHeight - 250)
	uiWidth := int32(400)
	uiHeight := int32(200)
	uiX := int32(ScreenWidth - uiWidth - 50)
	
	// Draw menu background
	rl.DrawRectangle(uiX, uiY, uiWidth, uiHeight, rl.NewColor(80, 20, 80, 240))
	rl.DrawRectangleLines(uiX, uiY, uiWidth, uiHeight, rl.White)
	
	// Draw title
	spell := r.game.CurrentSpells[r.game.SelectedSpell]
	titleText := fmt.Sprintf("Cast %s", spell.Name)
	rl.DrawText(titleText, uiX+20, uiY+10, 24, rl.White)
	
	// Draw targets based on spell type
	var targets []*combat.Person
	if spell.TargetType == game.TargetEnemy {
		targets = r.game.EnemyParty
	} else {
		targets = r.game.PlayerParty
	}
	
	for i, target := range targets {
		if target.Health > 0 || spell.Type == game.SpellHeal {
			optionY := uiY + 50 + int32(i*30)
			textColor := rl.LightGray
			
			if i == r.game.SelectedTarget {
				// Highlight selected target
				rl.DrawRectangle(uiX+10, optionY-5, uiWidth-20, 25, rl.DarkPurple)
				textColor = rl.Yellow
				
				// Highlight the target character
				r.highlightSpellTarget(target, spell.TargetType)
			}
			
			var targetText string
			if spell.Type == game.SpellHeal {
				targetText = fmt.Sprintf("%s (HP: %d/%d)", target.Name, target.Health, r.getMaxHealthForTarget(target))
			} else {
				targetText = fmt.Sprintf("%s (HP: %d)", target.Name, target.Health)
			}
			rl.DrawText(targetText, uiX+20, optionY, 20, textColor)
		}
	}
	
	// Draw instructions
	rl.DrawText("↑↓ Navigate  ENTER Cast  ESC Back", uiX+10, uiY+uiHeight-25, 16, rl.LightGray)
}

func (r *Renderer) highlightSpellTarget(target *combat.Person, targetType game.TargetType) {
	var targets []*combat.Person
	var baseX int32
	
	if targetType == game.TargetEnemy {
		targets = r.game.EnemyParty
		baseX = int32(ScreenWidth - 200)
	} else {
		targets = r.game.PlayerParty
		baseX = 100
	}
	
	// Find target position and draw highlight
	for i, char := range targets {
		if char == target {
			var baseY int32
			if targetType == game.TargetEnemy {
				baseY = int32(250 + i*120)
			} else {
				baseY = int32(200 + i*120)
			}
			
			charWidth := int32(100)
			charHeight := int32(150)
			
			// Draw pulsing highlight
			alpha := 0.5 + 0.3*float32(rl.GetTime()*4)
			if alpha > 1 {
				alpha = 1
			}
			
			highlightColor := rl.Purple
			if targetType == game.TargetAlly {
				highlightColor = rl.Green
			}
			
			rl.DrawRectangleLines(baseX-10, baseY-10, charWidth+20, charHeight+20, rl.Fade(highlightColor, alpha))
			break
		}
	}
}

func (r *Renderer) getMaxHealthForTarget(target *combat.Person) int {
	switch target.Name {
	case "Warrior":
		return 120
	case "Mage":
		return 80
	case "Rogue":
		return 90
	case "Orc":
		return 100
	case "Goblin":
		return 70
	default:
		return 100
	}
}

func (r *Renderer) drawVictoryScreen() {
	r.drawEndScreen("VICTORY!", rl.Green)
}

func (r *Renderer) drawDefeatScreen() {
	r.drawEndScreen("DEFEAT!", rl.Red)
}

func (r *Renderer) drawEndScreen(text string, color rl.Color) {
	centerX := int32(ScreenWidth / 2)
	centerY := int32(ScreenHeight / 2)

	textSize := int32(80)
	textWidth := rl.MeasureText(text, textSize)
	rl.DrawText(text, centerX-textWidth/2, centerY-100, textSize, color)

	instructions := "Press SPACE to Return to Menu"
	instrSize := int32(30)
	instrWidth := rl.MeasureText(instructions, instrSize)
	rl.DrawText(instructions, centerX-instrWidth/2, centerY, instrSize, rl.LightGray)
}