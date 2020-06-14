# development
This is my "spec sheet" for designing a roguelike. I'm trying to boil down as much as I can while still keeping the elements that make the game fun. The most important element is that everything should be interactive with eachother. For example, if the game has a can of grease, you should be able to apply it to your armor/weapons to help keep it from rusting (until it washes off) and protect against grab attacks. You can also apply it to the ground/stairs to make a trap for enemies (or yourself!) to slip and fall, apply it to a stuck door/container to open it in fewer turns, apply it to a friendly creature to grease their armor/weapons, apply it to . The can should also run out of grease and be able to be filled with water for other uses. The grease can also be lit on fire to give off a dim light. If the grease was applied to the floor, the floor can be lit on fire, etc. The can of grease example is from this video: https://invidio.us/watch?v=SjuTyJlgLJ8

# plans
I plan to implement these different "levels" of functionality, keeping the game as generic as possible for as long as possible, so that I can eventually split the codebase off into a library and a game that uses that library. That way, more games can easily be developed based on that library.
Update: looks like RogueBasin has a similar [article](http://roguebasin.roguelikedevelopment.org/index.php?title=How_to_Write_a_Roguelike_in_15_Steps)

# levels

## level 0: movement
### tiles
* empty: .
* wall: #

### creatures
* player: @, highlighted

### generation
* generate 1 48x16 level with each tile randomly empty or wall
* start the player on an empty square

### display
* display the map (all tiles visible)
* display the player, highlighted 

### player verbs
* move: numpad/arrow keys/p[;' (diagonals)/click screen



## level 1: levels
### tiles
* up stairs: <
* down stairs: >

### generation
* generate better looking levels
* up stairs and down stairs should be connected (flood fill)
* can have multiple up/down stairs?
* generate 32 levels
* player starts on upstairs of level 0
* first level has no up stairs
* last level has no down stairs

### display
* raycasting
* mask for explored regions
* mask for visible regions

### player verbs
* go up stairs: <
* go down stairs: >



## level 2: message system
### player verbs
* view message log: m (actually not implemented until tabs)
### messages
* running into a wall: "oof!"
* going up stairs: "you walk up the stairs"
* going down stairs: "you walk down the stairs"
* entering a new level: "this level you have entered looks completely new to you"
* note: must deal with multiple messages with "--more--"



## level 3: fighting

### creatures
* enemy: e

### generation
* generate enemes
* spawn enemies periodically

### player verbs
* attack: move, but towards an enemy
* kick: k, followed by move direction
	* another attack that deals less damage
	* why? for when more features are added it will become useful

### status
* HP
* death

### messages
* "an enemy approaches you"
* "you approach an enemy"
* "the enemy hits you"
* "you hit the enemy"
* "the enemy kicks you"
* "you kick the enemy"
* "the enemy dies"
* "you die!"


##level 4: items

### items
* weapon: )
* armor: [

### player verbs
* pick up: ,
* drop: d
* equip/unequip: e
* throw: t
* kick: k
	* kicking an item will propel it forward
	* kicking a creature will attack it, not adding your weapon damage

### status
* inventory sidebar
* damage output
* AC


##level 5: tabs
### display
* tabbed interface on the right

### tabs
* tutorial/instructions/help
* inventory
* list of player verbs
* message history

### player verbs
* view inventory: i
* view message history: m
* view help tab: ?
* view any tab: click








##level X:
### tiles
* 

### creatures
* 

### generation
* 

### display
* 

### player verbs
* 

### messages
* 

### status
* 

### tabs
* 







# things to add
* autowalk when screen clicked
* give more info when screen right clicked
* message system
* goal/win condition
* can of grease, duct tape
* item types
* better level gen
* food
* XP levels, player and enemy
* "magic" items
* special levels
* shops
* unidentified items
* dungeon features
* named items/creatures
* status effects
* different enemies
* raycasting, line of sight, unexplored
* hearing sounds from around the level
* walls, doors
* money
* fortune cookies :)
* autopickup and convenience features
* config file/options
* save files
* map branching, multiple up/down stairs
* traps
* peaceful creatures
* extended command menu (and quit)
* chatting
* stats
* classes
* races/species
* starting equipment
* inventory categories
* BUC
* -1, +0, +1 enchantment
* alignment? or factions
* gender
* light sources
* mix of randomly generated level and template structures
* colors
* falling down prone
* burdened
* bags
* prayer/phone home
* multishot
* rust, corrosion, burnt, etc.
* can of grease
* water
* teleportation
* minigames (like Sokoban)
* artifacts
* music
* play online
* "pop-up" windows
* farlook command
* alternate move methods
* check if screen is too small
* run in direction commands