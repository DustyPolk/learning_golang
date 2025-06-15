# Go Combat Arena - Final Fantasy Style Battle System

## How to Play

1. **Run the game**: `go run cmd/game/main.go`

2. **Main Menu**: Press SPACE to start a battle

3. **Party-Based Combat**:
   - **Your Party**: Warrior (tank), Mage (magic), Rogue (speed)
   - **Enemy Party**: Orc (strong), Goblin (fast)
   - Turn order based on Dodge stat (speed/agility)
   - Each character takes individual turns

4. **Battle Flow**:
   - When it's your character's turn, an action menu appears
   - Choose from: Attack, Magic, Item, or Run
   - Select your target from the enemy party
   - Watch the battle animations play out
   - Enemy AI automatically chooses actions and targets

## Controls

### Menu Navigation
- **↑/↓ Arrow Keys**: Navigate menus
- **ENTER or SPACE**: Select option
- **ESC**: Go back (in target selection)

### Battle States
- **Action Menu**: Choose what your character will do
- **Target Selection**: Choose which enemy to attack
- **Battle Execution**: Watch animations and effects

## Combat Features

### Attack Types
- **Melee**: +20% damage, character rushes to enemy with slash
- **Ranged**: Higher accuracy projectile attack
- **Special**: 2x damage explosion that cannot be dodged

### Party Roles
- **Warrior**: High HP (120), strong defense, reliable damage
- **Mage**: Low HP (80), high dodge, magical abilities
- **Rogue**: Medium HP (90), highest dodge, balanced stats

### Enemy AI
- **Smart Targeting**: Focuses on weakest party members
- **Adaptive Strategy**: Uses different attacks based on situation
- **Tactical Decisions**: Special attacks when desperate or facing strong enemies

## Visual Features

- **Party Layout**: Classic FF-style positioning (players left, enemies right)
- **Turn Order Display**: Shows upcoming 5 turns
- **Character Highlighting**: Current turn character glows with indicator
- **Target Selection**: Visual highlighting of selected enemy
- **Battle Status**: Real-time status messages
- **Health Bars**: Individual HP for each character
- **Animations**: Different effects for each attack type
- **Battle Log**: Scrolling combat information

## Battle UI

- **Turn Order Panel**: Top-right shows turn sequence
- **Action Menu**: Bottom-right during player turns
- **Target Selection**: Choose enemies with visual feedback
- **Status Display**: Shows current action state
- **Character Stats**: Hover info shows STR, DEF, DODGE

This creates a classic JRPG battle experience with modern visual effects!