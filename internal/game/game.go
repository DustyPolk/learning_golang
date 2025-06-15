package game

import (
	"fmt"
	"math/rand"
	"github.com/dustypolk/learning_golang/internal/combat"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	StateMenu GameState = iota
	StateBattle
	StateActionMenu
	StateTargetSelect
	StateSpellSelect
	StateSpellTarget
	StateVictory
	StateDefeat
)

type BattleMenuOption int

const (
	MenuAttack BattleMenuOption = iota
	MenuMagic
	MenuItem
	MenuRun
)

type Game struct {
	State              GameState
	PlayerParty        []*combat.Person
	EnemyParty         []*combat.Person
	TurnOrder          []*combat.Person
	CurrentTurnIndex   int
	CurrentCharacter   *combat.Person
	BattleLog          []string
	AnimationTimer     float32
	ActionQueue        []Action
	IsPaused           bool
	TurnDelay          float32
	NextTurnTime       float32
	
	// Menu system
	SelectedMenuOption BattleMenuOption
	SelectedTarget     int
	SelectedSpell      int
	IsSelectingTarget  bool
	SelectedAttackType AttackType
	CurrentSpells      []Spell
}

type Action struct {
	Type       ActionType
	AttackType AttackType
	Attacker   *combat.Person
	Defender   *combat.Person
	Damage     int
	IsCritical bool
	IsDodged   bool
	StartTime  float32
}

type ActionType int

const (
	ActionAttack ActionType = iota
	ActionDodge
	ActionHit
)

type AttackType int

const (
	AttackMelee AttackType = iota
	AttackRanged
	AttackSpecial
)

type SpellType int

const (
	SpellFireball SpellType = iota
	SpellHeal
	SpellLightning
	SpellShield
)

type Spell struct {
	Name        string
	Type        SpellType
	MPCost      int
	Power       int
	Description string
	TargetType  TargetType // Enemy, Ally, Self
}

type TargetType int

const (
	TargetEnemy TargetType = iota
	TargetAlly
	TargetSelf
)

func NewGame() *Game {
	return &Game{
		State:     StateMenu,
		BattleLog: make([]string, 0),
	}
}

func (g *Game) getSpellsForCharacter(character *combat.Person) []Spell {
	var spells []Spell
	
	switch character.Name {
	case "Mage":
		spells = []Spell{
			{Name: "Fireball", Type: SpellFireball, MPCost: 12, Power: 35, Description: "Fire damage", TargetType: TargetEnemy},
			{Name: "Heal", Type: SpellHeal, MPCost: 8, Power: 25, Description: "Restore HP", TargetType: TargetAlly},
			{Name: "Lightning", Type: SpellLightning, MPCost: 15, Power: 30, Description: "Electric damage", TargetType: TargetEnemy},
		}
	case "Warrior":
		spells = []Spell{
			{Name: "Heal", Type: SpellHeal, MPCost: 15, Power: 15, Description: "Basic healing", TargetType: TargetAlly},
		}
	case "Rogue":
		spells = []Spell{
			{Name: "Lightning", Type: SpellLightning, MPCost: 20, Power: 25, Description: "Quick strike", TargetType: TargetEnemy},
			{Name: "Heal", Type: SpellHeal, MPCost: 12, Power: 20, Description: "Quick heal", TargetType: TargetAlly},
		}
	case "Orc":
		spells = []Spell{
			{Name: "Fireball", Type: SpellFireball, MPCost: 10, Power: 20, Description: "Crude fire", TargetType: TargetEnemy},
		}
	case "Goblin":
		spells = []Spell{
			{Name: "Lightning", Type: SpellLightning, MPCost: 15, Power: 18, Description: "Zap", TargetType: TargetEnemy},
		}
	}
	
	// Filter spells by available MP
	var availableSpells []Spell
	for _, spell := range spells {
		if character.MP >= spell.MPCost {
			availableSpells = append(availableSpells, spell)
		}
	}
	
	return availableSpells
}

func (g *Game) StartBattle() {
	g.State = StateBattle
	
	// Create player party
	g.PlayerParty = []*combat.Person{
		{Name: "Warrior", Strength: 100, Health: 120, MP: 20, MaxMP: 20, Defense: 20, Armor: 15, Dodge: 10},
		{Name: "Mage", Strength: 60, Health: 80, MP: 50, MaxMP: 50, Defense: 5, Armor: 5, Dodge: 25},
		{Name: "Rogue", Strength: 80, Health: 90, MP: 30, MaxMP: 30, Defense: 10, Armor: 8, Dodge: 30},
	}
	
	// Create enemy party
	g.EnemyParty = []*combat.Person{
		{Name: "Orc", Strength: 85, Health: 100, MP: 15, MaxMP: 15, Defense: 15, Armor: 10, Dodge: 15},
		{Name: "Goblin", Strength: 65, Health: 70, MP: 25, MaxMP: 25, Defense: 8, Armor: 5, Dodge: 35},
	}
	
	g.BattleLog = []string{"Battle Start!"}
	g.setupTurnOrder()
}

func (g *Game) setupTurnOrder() {
	// Combine all living characters
	g.TurnOrder = make([]*combat.Person, 0)
	
	for _, char := range g.PlayerParty {
		if char.Health > 0 {
			g.TurnOrder = append(g.TurnOrder, char)
		}
	}
	
	for _, char := range g.EnemyParty {
		if char.Health > 0 {
			g.TurnOrder = append(g.TurnOrder, char)
		}
	}
	
	// Simple turn order based on dodge stat (agility)
	for i := 0; i < len(g.TurnOrder); i++ {
		for j := i + 1; j < len(g.TurnOrder); j++ {
			if g.TurnOrder[i].Dodge < g.TurnOrder[j].Dodge {
				g.TurnOrder[i], g.TurnOrder[j] = g.TurnOrder[j], g.TurnOrder[i]
			}
		}
	}
	
	g.CurrentTurnIndex = 0
	if len(g.TurnOrder) > 0 {
		g.CurrentCharacter = g.TurnOrder[0]
		g.AddToLog(g.CurrentCharacter.Name + " goes first!")
		
		// If it's a player character, show action menu
		if g.isPlayerCharacter(g.CurrentCharacter) {
			g.State = StateActionMenu
		}
	}
}

func (g *Game) isPlayerCharacter(character *combat.Person) bool {
	for _, char := range g.PlayerParty {
		if char == character {
			return true
		}
	}
	return false
}

func (g *Game) nextTurn() {
	g.CurrentTurnIndex++
	if g.CurrentTurnIndex >= len(g.TurnOrder) {
		g.CurrentTurnIndex = 0
	}
	
	// Skip dead characters
	for i := 0; i < len(g.TurnOrder); i++ {
		if g.TurnOrder[g.CurrentTurnIndex].Health > 0 {
			break
		}
		g.CurrentTurnIndex++
		if g.CurrentTurnIndex >= len(g.TurnOrder) {
			g.CurrentTurnIndex = 0
		}
	}
	
	g.CurrentCharacter = g.TurnOrder[g.CurrentTurnIndex]
	
	// Check for battle end
	if g.checkBattleEnd() {
		return
	}
	
	// Set up next turn
	if g.isPlayerCharacter(g.CurrentCharacter) {
		g.State = StateActionMenu
		g.SelectedMenuOption = MenuAttack
	} else {
		g.State = StateBattle
		g.NextTurnTime = g.AnimationTimer + 1.0
	}
}

func (g *Game) checkBattleEnd() bool {
	playerAlive := false
	enemyAlive := false
	
	for _, char := range g.PlayerParty {
		if char.Health > 0 {
			playerAlive = true
			break
		}
	}
	
	for _, char := range g.EnemyParty {
		if char.Health > 0 {
			enemyAlive = true
			break
		}
	}
	
	if !playerAlive {
		g.State = StateDefeat
		return true
	} else if !enemyAlive {
		g.State = StateVictory
		return true
	}
	
	return false
}

func (g *Game) AddToLog(message string) {
	g.BattleLog = append(g.BattleLog, message)
	if len(g.BattleLog) > 10 {
		g.BattleLog = g.BattleLog[1:]
	}
}

func (g *Game) Update(deltaTime float32) {
	g.AnimationTimer += deltaTime

	switch g.State {
	case StateBattle:
		g.updateBattle(deltaTime)
	}
}

func (g *Game) updateBattle(deltaTime float32) {
	if len(g.ActionQueue) > 0 {
		action := &g.ActionQueue[0]
		animDuration := float32(0.5)
		
		// Different durations for different attacks
		switch action.AttackType {
		case AttackMelee:
			animDuration = 0.8
		case AttackRanged:
			animDuration = 0.5
		case AttackSpecial:
			animDuration = 1.0
		}
		
		if g.AnimationTimer-action.StartTime > animDuration {
			g.ActionQueue = g.ActionQueue[1:]
			// After animation ends, move to next turn
			if len(g.ActionQueue) == 0 {
				g.nextTurn()
			}
		}
	}

	// Handle turn advancement after delay (when no animations)
	if g.State == StateBattle && len(g.ActionQueue) == 0 && g.AnimationTimer >= g.NextTurnTime {
		if !g.isPlayerCharacter(g.CurrentCharacter) {
			// Enemy turn
			if g.CurrentCharacter.Health > 0 {
				g.enemyAI()
			} else {
				g.nextTurn()
			}
		}
	}
}

func (g *Game) HandleInput() {
	switch g.State {
	case StateMenu:
		if rl.IsKeyPressed(rl.KeySpace) {
			g.StartBattle()
		}
	case StateActionMenu:
		// Navigate action menu
		if rl.IsKeyPressed(rl.KeyUp) {
			if g.SelectedMenuOption > 0 {
				g.SelectedMenuOption--
			}
		} else if rl.IsKeyPressed(rl.KeyDown) {
			if g.SelectedMenuOption < MenuRun {
				g.SelectedMenuOption++
			}
		} else if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
			g.selectMenuOption()
		}
	case StateTargetSelect:
		// Navigate target selection
		if rl.IsKeyPressed(rl.KeyUp) {
			if g.SelectedTarget > 0 {
				g.SelectedTarget--
			}
		} else if rl.IsKeyPressed(rl.KeyDown) {
			maxTargets := len(g.EnemyParty) - 1
			if g.SelectedTarget < maxTargets {
				g.SelectedTarget++
			}
		} else if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
			g.selectTarget()
		} else if rl.IsKeyPressed(rl.KeyEscape) {
			g.State = StateActionMenu
		}
	case StateSpellSelect:
		// Navigate spell selection
		if rl.IsKeyPressed(rl.KeyUp) {
			if g.SelectedSpell > 0 {
				g.SelectedSpell--
			}
		} else if rl.IsKeyPressed(rl.KeyDown) {
			maxSpells := len(g.CurrentSpells) - 1
			if g.SelectedSpell < maxSpells {
				g.SelectedSpell++
			}
		} else if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
			g.selectSpell()
		} else if rl.IsKeyPressed(rl.KeyEscape) {
			g.State = StateActionMenu
		}
	case StateSpellTarget:
		// Navigate spell target selection
		if rl.IsKeyPressed(rl.KeyUp) {
			if g.SelectedTarget > 0 {
				g.SelectedTarget--
			}
		} else if rl.IsKeyPressed(rl.KeyDown) {
			spell := g.CurrentSpells[g.SelectedSpell]
			var maxTargets int
			if spell.TargetType == TargetEnemy {
				maxTargets = len(g.EnemyParty) - 1
			} else {
				maxTargets = len(g.PlayerParty) - 1
			}
			if g.SelectedTarget < maxTargets {
				g.SelectedTarget++
			}
		} else if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
			g.castSpell()
		} else if rl.IsKeyPressed(rl.KeyEscape) {
			g.State = StateSpellSelect
		}
	case StateBattle:
		// Enemy turn processing
		if !g.isPlayerCharacter(g.CurrentCharacter) && g.AnimationTimer >= g.NextTurnTime && len(g.ActionQueue) == 0 {
			if g.CurrentCharacter.Health > 0 {
				g.enemyAI()
			}
		}
	case StateVictory, StateDefeat:
		if rl.IsKeyPressed(rl.KeySpace) {
			g.State = StateMenu
		}
	}
}

func (g *Game) selectMenuOption() {
	switch g.SelectedMenuOption {
	case MenuAttack:
		g.State = StateTargetSelect
		g.SelectedTarget = 0
		g.SelectedAttackType = AttackMelee // Default to melee for now
	case MenuMagic:
		g.CurrentSpells = g.getSpellsForCharacter(g.CurrentCharacter)
		if len(g.CurrentSpells) == 0 {
			g.AddToLog(g.CurrentCharacter.Name + " has no spells available!")
			g.State = StateActionMenu
		} else {
			g.State = StateSpellSelect
			g.SelectedSpell = 0
		}
	case MenuItem:
		g.AddToLog(g.CurrentCharacter.Name + " uses an item!")
		// TODO: Implement item system
		g.endTurn()
	case MenuRun:
		g.AddToLog("Tried to run away!")
		g.endTurn()
	}
}

func (g *Game) selectSpell() {
	if g.SelectedSpell < len(g.CurrentSpells) {
		g.State = StateSpellTarget
		g.SelectedTarget = 0
	}
}

func (g *Game) castSpell() {
	if g.SelectedSpell < len(g.CurrentSpells) {
		spell := g.CurrentSpells[g.SelectedSpell]
		
		// Check MP cost
		if g.CurrentCharacter.MP < spell.MPCost {
			g.AddToLog(g.CurrentCharacter.Name + " doesn't have enough MP!")
			g.State = StateActionMenu
			return
		}
		
		// Deduct MP
		g.CurrentCharacter.MP -= spell.MPCost
		
		// Get target
		var target *combat.Person
		if spell.TargetType == TargetEnemy {
			if g.SelectedTarget < len(g.EnemyParty) {
				target = g.EnemyParty[g.SelectedTarget]
			}
		} else {
			if g.SelectedTarget < len(g.PlayerParty) {
				target = g.PlayerParty[g.SelectedTarget]
			}
		}
		
		if target == nil {
			g.AddToLog("Invalid target!")
			g.State = StateActionMenu
			return
		}
		
		// Apply spell effect
		g.applySpellEffect(spell, target)
		g.State = StateBattle
	}
}

func (g *Game) selectTarget() {
	if g.SelectedTarget < len(g.EnemyParty) {
		target := g.EnemyParty[g.SelectedTarget]
		if target.Health > 0 {
			g.PerformAttackWithType(g.CurrentCharacter, target, g.SelectedAttackType)
			g.State = StateBattle
		}
	}
}

func (g *Game) endTurn() {
	g.State = StateBattle
	g.NextTurnTime = g.AnimationTimer + 1.0
}

func (g *Game) PerformAttackWithType(attacker, defender *combat.Person, attackType AttackType) {
	var damage int
	var isCritical, isDodged bool
	var attackName string
	
	switch attackType {
	case AttackMelee:
		// Melee: Higher damage, normal crit chance
		damage, isCritical, isDodged = combat.CalculateDamageWithInfo(attacker, defender)
		damage = int(float64(damage) * 1.2) // 20% more damage
		attackName = "melee strikes"
	case AttackRanged:
		// Ranged: Normal damage, higher accuracy
		damage, isCritical, isDodged = combat.CalculateDamageWithInfo(attacker, defender)
		// Reduce dodge chance for ranged attacks
		if isDodged && defender.Dodge > 10 {
			isDodged = false
			damage = int(float64(damage) * 0.8) // But slightly less damage
		}
		attackName = "shoots"
	case AttackSpecial:
		// Special: High damage, always hits, but has cooldown
		damage, isCritical, _ = combat.CalculateDamageWithInfo(attacker, defender)
		damage = int(float64(damage) * 2.0) // Double damage
		isDodged = false // Special attacks can't be dodged
		attackName = "unleashes SPECIAL on"
	}
	
	action := Action{
		Type:       ActionAttack,
		AttackType: attackType,
		Attacker:   attacker,
		Defender:   defender,
		Damage:     damage,
		IsCritical: isCritical,
		IsDodged:   isDodged,
		StartTime:  g.AnimationTimer,
	}
	g.ActionQueue = append(g.ActionQueue, action)

	if isDodged {
		g.AddToLog(defender.Name + " dodged the attack!")
	} else if isCritical {
		g.AddToLog(fmt.Sprintf("%s %s %s for %d CRITICAL damage!", attacker.Name, attackName, defender.Name, damage))
	} else {
		g.AddToLog(fmt.Sprintf("%s %s %s for %d damage!", attacker.Name, attackName, defender.Name, damage))
	}

	if !isDodged {
		defender.Health -= damage
	}

	// Set delay for next turn
	g.NextTurnTime = g.AnimationTimer + 1.5
}

func (g *Game) PerformAttack(attacker, defender *combat.Person) {
	// Default to ranged attack for backwards compatibility
	g.PerformAttackWithType(attacker, defender, AttackRanged)
}

func (g *Game) enemyAI() {
	// Choose target from player party
	var target *combat.Person
	aliveTargets := make([]*combat.Person, 0)
	
	for _, char := range g.PlayerParty {
		if char.Health > 0 {
			aliveTargets = append(aliveTargets, char)
		}
	}
	
	if len(aliveTargets) == 0 {
		return
	}
	
	// Simple AI: target lowest health player
	target = aliveTargets[0]
	for _, char := range aliveTargets {
		if char.Health < target.Health {
			target = char
		}
	}
	
	// AI decision making based on game state
	enemyMaxHealth := 100
	if g.CurrentCharacter.Name == "Orc" {
		enemyMaxHealth = 100
	} else if g.CurrentCharacter.Name == "Goblin" {
		enemyMaxHealth = 70
	}
	
	healthPercent := float32(g.CurrentCharacter.Health) / float32(enemyMaxHealth)
	targetMaxHealth := 120
	if target.Name == "Mage" {
		targetMaxHealth = 80
	} else if target.Name == "Rogue" {
		targetMaxHealth = 90
	}
	targetHealthPercent := float32(target.Health) / float32(targetMaxHealth)
	
	var attackChoice AttackType
	
	aiRoll := rand.Float32()
	
	if healthPercent < 0.3 && aiRoll < 0.6 {
		attackChoice = AttackSpecial
		fmt.Printf("AI: %s desperate, using SPECIAL!\n", g.CurrentCharacter.Name)
	} else if targetHealthPercent < 0.3 && aiRoll < 0.7 {
		attackChoice = AttackMelee
		fmt.Printf("AI: %s targeting weak %s with melee!\n", g.CurrentCharacter.Name, target.Name)
	} else {
		choice := rand.Intn(100)
		if choice < 50 {
			attackChoice = AttackRanged
			fmt.Printf("AI: %s using ranged on %s\n", g.CurrentCharacter.Name, target.Name)
		} else if choice < 80 {
			attackChoice = AttackMelee
			fmt.Printf("AI: %s using melee on %s\n", g.CurrentCharacter.Name, target.Name)
		} else {
			attackChoice = AttackSpecial
			fmt.Printf("AI: %s surprise SPECIAL on %s!\n", g.CurrentCharacter.Name, target.Name)
		}
	}
	
	g.PerformAttackWithType(g.CurrentCharacter, target, attackChoice)
}

func (g *Game) applySpellEffect(spell Spell, target *combat.Person) {
	switch spell.Type {
	case SpellFireball, SpellLightning:
		// Damage spells
		damage := spell.Power + rand.Intn(10) - 5 // Some variance
		if damage < 1 {
			damage = 1
		}
		
		// Check for dodge (spells are harder to dodge)
		dodgeRoll := rand.Intn(100)
		if dodgeRoll < target.Dodge/2 { // Half dodge chance vs magic
			g.AddToLog(target.Name + " dodged the " + spell.Name + "!")
		} else {
			target.Health -= damage
			if target.Health < 0 {
				target.Health = 0
			}
			
			spellEffect := "fire"
			if spell.Type == SpellLightning {
				spellEffect = "lightning"
			}
			
			g.AddToLog(fmt.Sprintf("%s casts %s on %s for %d %s damage!", 
				g.CurrentCharacter.Name, spell.Name, target.Name, damage, spellEffect))
		}
		
		// Create magic action for animation
		action := Action{
			Type:       ActionAttack,
			AttackType: AttackSpecial, // Use special for magic effects
			Attacker:   g.CurrentCharacter,
			Defender:   target,
			Damage:     damage,
			IsCritical: false,
			IsDodged:   dodgeRoll < target.Dodge/2,
			StartTime:  g.AnimationTimer,
		}
		g.ActionQueue = append(g.ActionQueue, action)
		
	case SpellHeal:
		// Healing spell
		healing := spell.Power + rand.Intn(8) - 2 // Some variance
		if healing < 1 {
			healing = 1
		}
		
		maxHealth := g.getMaxHealthForCharacter(target)
		target.Health += healing
		if target.Health > maxHealth {
			target.Health = maxHealth
		}
		
		g.AddToLog(fmt.Sprintf("%s casts Heal on %s, restoring %d HP!", 
			g.CurrentCharacter.Name, target.Name, healing))
			
	case SpellShield:
		// Buff spell (placeholder for now)
		g.AddToLog(fmt.Sprintf("%s casts Shield on %s!", 
			g.CurrentCharacter.Name, target.Name))
	}
	
	// Set delay for next turn
	g.NextTurnTime = g.AnimationTimer + 1.5
}

func (g *Game) getMaxHealthForCharacter(character *combat.Person) int {
	switch character.Name {
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