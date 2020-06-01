package main

//Position - the x, y, and z (depth) of a location in the dungeon
type Position struct {
	x int
	y int
	z int
}

//Point - the x, and y location on a level
type Point struct {
	x int
	y int
}

//Dungeon - the layout of the entire dungeon, including all the levels and metadata
type Dungeon struct {
	levels [depth]*Level
}

//Level - one layer of the dungeon, made of many tiles
type Level struct {
	tiles [height][width]*Tile
}

//Tile - one square on a level
type Tile struct {
	char int32 //the character that is displayed
	// isSolid     bool  //are you prevented from walking through it?
	// blocksLight bool  //are you prevented from seeing through it?
	// isRoom      bool  //does it define the walls of a room?
	// isCorner    bool  //is it the corner of a room?
}

//GetTile - retrieve a tile from the dungeon using an x, y, z position
func (d *Dungeon) GetTile(p Position) *Tile {
	return d.levels[p.z].tiles[p.y][p.x]
}

//SetTile - set a tile in the dungeon using an x, y, z position
func (d *Dungeon) SetTile(p Position, t *Tile) {
	d.levels[p.z].tiles[p.y][p.x] = t
}

//GetTile - retrieve a tile from a level using an x, y, z position
func (l *Level) GetTile(p Point) *Tile {
	return l.tiles[p.y][p.x]
}

//SetTile - set a tile in a level using an x, y, z position
func (l *Level) SetTile(p Point, t *Tile) {
	l.tiles[p.y][p.x] = t
}

//GetChar - retrieve a char from a level using an x, y position
func (l *Level) GetChar(p Point) int32 {
	return l.tiles[p.y][p.x].char
}

//IsSolid - returns whether or not a creature can move through the tile
func (t *Tile) IsSolid() bool {
	c := t.char
	return c == ' ' || c == '#' || c == '─' || c == '│' || t.IsCorner()
}

//IsCorner - is the tile the corner of a room border?
func (t *Tile) IsCorner() bool {
	c := t.char
	return c == '├' || c == '┤' || c == '┬' || c == '┴' || c == '┼' || c == '┌' || c == '└' || c == '┐' || c == '┘'
}

//IsRoom - is the tile part of a room border, including doors?
func (t *Tile) IsRoom() bool {
	c := t.char
	return c == '─' || c == '│' || c == '*' || t.IsCorner()
}

//IsRoomFloor - is the tile one that is normally generated inside a room?
func (t *Tile) IsRoomFloor() bool {
	c := t.char
	return c == '·' || c == '>' || c == '<' //|| c == '*'
}

//BlocksLight - returns whether or not a creature can see through the tile
func (t *Tile) BlocksLight() bool {
	return t.IsSolid()
}

//SetChar - set a char in a level using an x, y, z position
func (l *Level) SetChar(p Point, c int32) {
	l.tiles[p.y][p.x] = &Tile{c}
}

//GetChar - retrieve a char from the dungeon using an x, y, z position
func (d *Dungeon) GetChar(p Position) int32 {
	return d.GetTile(p).char
}

//SetChar - set a char in a level using an x, y, z position
func (d *Dungeon) SetChar(p Position, c int32) {
	d.SetTile(p, &Tile{c})
}

//GetLevel - get one level from the dungeon
func (d *Dungeon) GetLevel(z int) *Level {
	return d.levels[z]
}

//NewLevel - create a new empty level
func NewLevel() *Level {
	var tiles [height][width]*Tile
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			tiles[y][x] = &Tile{' '}
		}
	}
	return &Level{tiles}
}

//NewDungeon - create a new dungeon
func NewDungeon() *Dungeon {
	var levels [depth]*Level
	for z := 0; z < depth; z++ {
		levels[z] = NewLevel()
	}
	return &Dungeon{levels}
}

//FindChar - find the location of a particular character in the level
func (l *Level) FindChar(c int32) *Point {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if l.GetChar(Point{x, y}) == c {
				return &Point{x, y}
			}
		}
	}
	return nil
}

//DistSquaredTo - return the square of the distance to another point
func (p1 *Point) DistSquaredTo(p2 *Point) int {
	return (p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)
}
