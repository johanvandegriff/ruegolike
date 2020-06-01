package main

import (
	"math/rand"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func randRangeInclusive(min, max int) int {
	return rand.Intn(max-min+1) + min
}

type room struct {
	x, y, w, h int
}

//https://www.geeksforgeeks.org/find-two-rectangles-overlap/
func doRoomsOverlap(rm1, rm2 room) bool {

	// If one rectangle is on left side of other
	if rm1.x >= rm2.x+rm2.w-1 || rm2.x >= rm1.x+rm1.w-1 {
		return false
	}

	// If one rectangle is above other
	if rm1.y >= rm2.y+rm2.h-1 || rm2.y >= rm1.y+rm1.h-1 {
		return false
	}

	return true
}

//TODO maybe reject levels with too much empty space
//TODO different distributions of number of rooms to size of rooms
//	either lots of small rooms, a few big rooms, or in between
func genRoomLevel(level *Level) {
	for yi := 0; yi < height; yi++ {
		for xi := 0; xi < width; xi++ {
			level.SetChar(Point{xi, yi}, '#')
		}
	}

	numRooms := randRangeInclusive(4, 8) //random number of rooms
	rooms := make([]room, numRooms)
	for k := 0; k < numRooms; k++ {
		x, y, w, h := rand.Intn(width), rand.Intn(height), randRangeInclusive(4, 12), randRangeInclusive(4, 8)
		// x, y, w, h := rand.Intn(width), rand.Intn(height), randRangeInclusive(4, 4), randRangeInclusive(4, 4)

		//center the x and y with the width and height to get the same number of rooms hitting the top and
		//left walls as the bottom and right walls
		x -= w / 2
		y -= h / 2

		//limit the x and y positions to prevent array out of bounds
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		if x+w > width {
			x = width - w
		}
		if y+h > height {
			y = height - h
		}

		newRoom := room{x, y, w, h}

		//store the room
		rooms[k] = newRoom

		//check if the room intersects with any other rooms
		for k2 := 0; k2 < k; k2++ {
			if doRoomsOverlap(newRoom, rooms[k2]) {
				k--
				break
			}
		}

	}

	for k := 0; k < numRooms; k++ {
		x, y, w, h := rooms[k].x, rooms[k].y, rooms[k].w, rooms[k].h

		//top left corner
		addBoxArt(level, y, x, '┌')

		//bottom left corner
		addBoxArt(level, y+h-1, x, '└')

		//top right corner
		addBoxArt(level, y, x+w-1, '┐')

		//bottom right corner
		addBoxArt(level, y+h-1, x+w-1, '┘')

		//top and bottom walls
		for i := 1; i < w-1; i++ {
			addBoxArt(level, y, x+i, '─')
			addBoxArt(level, y+h-1, x+i, '─')
		}

		//left and right walls
		for j := 1; j < h-1; j++ {
			addBoxArt(level, y+j, x, '│')
			addBoxArt(level, y+j, x+w-1, '│')
		}

		//floor
		for i := 1; i < w-1; i++ {
			for j := 1; j < h-1; j++ {
				addBoxArt(level, y+j, x+i, '·')
			}
		}
	}

	start := rand.Intn(numRooms)
	end := rand.Intn(numRooms - 1)
	if start == end {
		end++
	}
	extra := 2 //TODO make sure extra corridors are not redundant, and use unused space on the outside of the map?
	//	also maybe add some non-winding corridors
	//	could use "concentric" rectangles starting from the outside and find intersection with rooms
	for {
		// for q := 0; q < 4 && !allConnected; q++ {
		tryDrawCorridor(start, end, rooms, level)

		//find the unconnected rooms
		connectedTo := make([]int, numRooms)
		for i := 0; i < numRooms; i++ {
			connectedTo[i] = i //each room starts off connected to itself
		}
		//now find the lowest room number that each room is connected to
		for i := 0; i < numRooms; i++ {
			if connectedTo[i] == i {
				//use flood fill to find what rooms each room is connected to
				mask := floodFill(rooms[i].x+1, rooms[i].y+1, level)
				for j := i + 1; j < numRooms; j++ {
					x, y := rooms[j].x+1, rooms[j].y+1
					if mask[y][x] {
						connectedTo[j] = i
					}
				}
			}
		}

		//TODO allow 1 unconnected room?
		numConnected := 0
		for i := 0; i < numRooms; i++ {
			if connectedTo[i] != 0 {
				numConnected++
				// break
			}
		}

		if numConnected == 0 {
			if extra == 0 {
				break
			}
			extra--
			start = rand.Intn(numRooms)
			end = rand.Intn(numRooms)
		} else {

			//pick what 2 rooms to connect next
			start = rand.Intn(numRooms)
			end = rand.Intn(numRooms)
			for start == end || connectedTo[start] == connectedTo[end] {
				start = rand.Intn(numRooms)
				end = rand.Intn(numRooms)
			}
		}
	}
}

func tryDrawCorridor(i1, i2 int, rooms []room, level *Level) bool {
	if i1 == i2 {
		return false
	}
	r1, r2 := rooms[i1], rooms[i2]
	//pick x,y locations inside the rooms not including the walls
	startX := randRangeInclusive(r1.x+1, r1.x+r1.w-2)
	startY := randRangeInclusive(r1.y+1, r1.y+r1.h-2)
	endX := randRangeInclusive(r2.x+1, r2.x+r2.w-2)
	endY := randRangeInclusive(r2.y+1, r2.y+r2.h-2)

	dx, dy := 0, 0
	x, y := startX, startY
	goX := true //slightly favor going horizontal first
	nextSame := false

	points := make([]Point, 0)

	for {
		if x < endX {
			dx = 1
		} else if x > endX {
			dx = -1
		} else {
			dx = 0
		}

		if y < endY {
			dy = 1
		} else if y > endY {
			dy = -1
		} else {
			dy = 0
		}

		if dx == 0 && dy == 0 {
			break
		}
		if dx == 0 {
			goX = false
		} else if dy == 0 {
			goX = true
		} else {
			//change direction 25% of the time
			if !nextSame && rand.Intn(100) < 25 {
				goX = !goX
			}
		}

		if goX {
			x += dx
		} else {
			y += dy
		}

		//stop when it hits another room other than intended. if not close enough, abort. if close enough, keep it
		if !level.GetTile(Point{x, y}).IsSolid() && nextSame {
			if (x-endX)*(x-endX)+(y-endY)*(y-endY) <= 8*8 {
				break
			} else {
				return false
			}
		}

		//stop when it hits a corner of a room
		if level.GetTile(Point{x, y}).IsCorner() {
			return false
		}

		nextSame = false
		if level.GetChar(Point{x, y}) == '─' {
			//prevent doors next to each other
			if level.GetChar(Point{x - 1, y}) == '*' || level.GetChar(Point{x + 1, y}) == '*' {
				return false
			}
			nextSame = true
		}
		if level.GetChar(Point{x, y}) == '│' {
			//prevent doors next to each other
			if level.GetChar(Point{x, y - 1}) == '*' || level.GetChar(Point{x, y + 1}) == '*' {
				return false
			}
			nextSame = true
		}

		points = append(points, Point{x, y})
	}

	for _, pt := range points {
		c := level.GetChar(pt)
		if c == '─' || c == '│' {
			level.SetChar(pt, '*')
		} else if c == '#' {
			level.SetChar(pt, ':')
		}
	}

	//TODO: option for "diagonal" tunnel (alternating x/y)
	return true
}

// https://unicode-search.net/unicode-namesearch.pl?term=BOX%20DRAWINGS
func addBoxArt(level *Level, y, x int, new int32) {
	old := level.GetChar(Point{x, y})
	combined := new

	//TODO replace this madness with contextual rendering
	switch old {
	case '┌':
		switch new {
		case '┌':
			combined = '┌'
		case '└':
			combined = '├'
		case '┐':
			combined = '┬'
		case '┘':
			combined = '┼'
		case '─':
			combined = '┬'
		case '│':
			combined = '├'
		case '├':
			combined = '├'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┬'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '└':
		switch new {
		case '┌':
			combined = '├'
		case '└':
			combined = '└'
		case '┐':
			combined = '┼'
		case '┘':
			combined = '┴'
		case '─':
			combined = '┴'
		case '│':
			combined = '├'
		case '├':
			combined = '├'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┼'
		case '┴':
			combined = '┴'
		case '┼':
			combined = '┼'
		}
	case '┐':
		switch new {
		case '┌':
			combined = '┬'
		case '└':
			combined = '┼'
		case '┐':
			combined = '┐'
		case '┘':
			combined = '┤'
		case '─':
			combined = '┬'
		case '│':
			combined = '┤'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┤'
		case '┬':
			combined = '┬'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '┘':
		switch new {
		case '┌':
			combined = '┼'
		case '└':
			combined = '┴'
		case '┐':
			combined = '┤'
		case '┘':
			combined = '┘'
		case '─':
			combined = '┴'
		case '│':
			combined = '┤'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┤'
		case '┬':
			combined = '┼'
		case '┴':
			combined = '┴'
		case '┼':
			combined = '┼'
		}
	case '─':
		switch new {
		case '┌':
			combined = '┬'
		case '└':
			combined = '┴'
		case '┐':
			combined = '┬'
		case '┘':
			combined = '┴'
		case '─':
			combined = '─'
		case '│':
			combined = '┼'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┬'
		case '┴':
			combined = '┴'
		case '┼':
			combined = '┼'
		}
	case '│':
		switch new {
		case '┌':
			combined = '├'
		case '└':
			combined = '├'
		case '┐':
			combined = '┤'
		case '┘':
			combined = '┤'
		case '─':
			combined = '┼'
		case '│':
			combined = '│'
		case '├':
			combined = '├'
		case '┤':
			combined = '┤'
		case '┬':
			combined = '┼'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '├':
		switch new {
		case '┌':
			combined = '├'
		case '└':
			combined = '├'
		case '┐':
			combined = '┼'
		case '┘':
			combined = '┼'
		case '─':
			combined = '┼'
		case '│':
			combined = '├'
		case '├':
			combined = '├'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┼'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '┤':
		switch new {
		case '┌':
			combined = '┼'
		case '└':
			combined = '┼'
		case '┐':
			combined = '┤'
		case '┘':
			combined = '┤'
		case '─':
			combined = '┼'
		case '│':
			combined = '┤'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┤'
		case '┬':
			combined = '┼'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '┬':
		switch new {
		case '┌':
			combined = '┬'
		case '└':
			combined = '┼'
		case '┐':
			combined = '┬'
		case '┘':
			combined = '┼'
		case '─':
			combined = '┬'
		case '│':
			combined = '┼'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┬'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '┴':
		switch new {
		case '┌':
			combined = '┼'
		case '└':
			combined = '┴'
		case '┐':
			combined = '┼'
		case '┘':
			combined = '┴'
		case '─':
			combined = '┴'
		case '│':
			combined = '┼'
		case '├':
			combined = '┼'
		case '┤':
			combined = '┼'
		case '┬':
			combined = '┴'
		case '┴':
			combined = '┼'
		case '┼':
			combined = '┼'
		}
	case '┼':
		combined = '┼'
	}
	level.SetChar(Point{x, y}, combined)
}

// http://roguebasin.roguelikedevelopment.org/index.php?title=Cellular_Automata_Method_for_Generating_Random_Cave-Like_Levels
func genCaveLevel(level *Level) {
	//48 16  40  5 1 4  5 0 3
	fillprob := 40

	//                 r1,r2,reps r1,r2,reps
	gens := [...][3]int{{5, 1, 4}, {5, 0, 3}}

	// var grid2 [height][width]int32
	level2 := NewLevel()

	for yi := 0; yi < height; yi++ {
		for xi := 0; xi < width; xi++ {
			if rand.Intn(100) < fillprob {
				level.SetChar(Point{xi, yi}, '#') //wall, 40%
			} else {
				level.SetChar(Point{xi, yi}, '·') //empty, 60%
			}
			level2.SetChar(Point{xi, yi}, '#')
		}
	}
	//border around the edge
	for yi := 0; yi < height; yi++ {
		level.SetChar(Point{0, yi}, '#')
		level.SetChar(Point{width - 1, yi}, '#')
	}
	for xi := 0; xi < width; xi++ {
		level.SetChar(Point{xi, 0}, '#')
		level.SetChar(Point{xi, height - 1}, '#')
	}

	for i := 0; i < len(gens); i++ {
		gen := gens[i]
		r1Cutoff := gen[0]
		r2Cutoff := gen[1]
		reps := gen[2]

		for j := 0; j < reps; j++ {

			for yi := 1; yi < height-1; yi++ {
				for xi := 1; xi < width-1; xi++ {
					adjCountR1, adjCountR2 := 0, 0

					for ii := -1; ii <= 1; ii++ {
						for jj := -1; jj <= 1; jj++ {
							if level.GetChar(Point{xi + jj, yi + ii}) != '·' {
								adjCountR1++
							}
						}
					}
					for ii := yi - 2; ii <= yi+2; ii++ {
						for jj := xi - 2; jj <= xi+2; jj++ {
							if abs(ii-yi) == 2 && abs(jj-xi) == 2 {
								continue
							}
							if ii < 0 || jj < 0 || ii >= height || jj >= width {
								continue
							}
							if level.GetChar(Point{jj, ii}) != '·' {
								adjCountR2++
							}
						}
					}
					if adjCountR1 >= r1Cutoff || adjCountR2 <= r2Cutoff {
						level2.SetChar(Point{xi, yi}, '#')
					} else {
						level2.SetChar(Point{xi, yi}, '·')
					}
				}
			}
			for yi := 1; yi < height-1; yi++ {
				for xi := 1; xi < width-1; xi++ {
					level.SetTile(Point{xi, yi}, level2.GetTile(Point{xi, yi}))
				}
			}
		}
	}

}

func floodFill(x, y int, level *Level) [height][width]bool {
	var mask [height][width]bool
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			floodFillAux(x+dx, y+dy, level, &mask)
		}
	}
	// floodFillAux(x, y, level, &mask)
	return mask
}

func floodFillAux(x, y int, level *Level, mask *[height][width]bool) {
	if isXYInRange(x, y) && !level.GetTile(Point{x, y}).IsSolid() && mask[y][x] == false {
		mask[y][x] = true
		// level[y][x] = '~'
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				floodFillAux(x+dx, y+dy, level, mask)
			}
		}
	}
}

const minStairDist = 8 //TODO experiment with this to avoid infinite loops

func tryToAddStairs(level *Level) bool {
	// var upStairX, upStairY, downStairX, downStairY int
	var upStair, downStair Point

	//find a place for the up stairs
	for i := 0; ; i++ {
		if i > 9999 {
			return false
		}
		upStair = Point{rand.Intn(width), rand.Intn(height)}
		if level.GetChar(Point{upStair.x, upStair.y}) == '·' {
			break
		}
	}

	mask := floodFill(upStair.x, upStair.y, level)
	for i := 0; ; i++ {
		if i > 9999 {
			return false
		}
		downStair = Point{rand.Intn(width), rand.Intn(height)}
		if mask[downStair.y][downStair.x] && level.GetChar(downStair) == '·' &&
			upStair.DistSquaredTo(&downStair) >= minStairDist*minStairDist {
			break
		}
	}
	level.SetChar(upStair, '>')
	level.SetChar(downStair, '<')
	return true

	// if z > 0 {
	// 	if z == 1 {
	// 		for i := 0; ; i++ {
	// 			if i > 9999 {
	// 				return false, 0, 0, 0, 0
	// 			}
	// 			stairX = rand.Intn(width)
	// 			stairY = rand.Intn(height)
	// 			if dungeon.GetChar(Position{stairX, stairY, z - 1}) == '·' {
	// 				break
	// 			}
	// 		}
	// 		playerX, playerY = stairX, stairY
	// 	}
	// 	mask := floodFill(stairX, stairY, dungeon.GetLevel(z-1))
	// 	oldStairX, oldStairY := stairX, stairY
	// 	for i := 0; ; i++ {
	// 		if i > 9999 {
	// 			return false, 0, 0, 0, 0
	// 		}
	// 		stairX = rand.Intn(width)
	// 		stairY = rand.Intn(height)
	// 		if mask[stairY][stairX] && dungeon.GetChar(Position{stairX, stairY, z - 1}) == '·' && dungeon.GetChar(Position{stairX, stairY, z}) == '·' &&
	// 			(stairX-oldStairX)*(stairX-oldStairX)+(stairY-oldStairY)*(stairY-oldStairY) >= minStairDist*minStairDist {
	// 			break
	// 		}
	// 	}
	// 	dungeon.SetChar(Position{stairX, stairY, z - 1}, '>')
	// 	dungeon.SetChar(Position{stairX, stairY, z}, '<')
	// }
	// return true, stairX, stairY, playerX, playerY
}

//Generate - generate all the levels in the game
func Generate() (*Dungeon, [depth][height][width]bool, Position) {
	// var levels [depth][height][width]int32
	dungeon := NewDungeon()
	var explored [depth][height][width]bool

	//simple terrain generation
	for z := 0; z < depth; z++ {
		roomType := rand.Intn(100) < 20
		succeeded := false

		for !succeeded {
			if roomType {
				genCaveLevel(dungeon.GetLevel(z))
			} else {
				genRoomLevel(dungeon.GetLevel(z))
			}
			succeeded = tryToAddStairs(dungeon.GetLevel(z))
		}
	}

	//replace the up stairs on the first level with a floor, and put the player there
	playerPos := dungeon.GetLevel(0).FindChar('<')
	dungeon.GetLevel(0).SetChar(*playerPos, '.')

	//remove the down stairs from the last level
	dungeon.GetLevel(depth-1).SetChar(*dungeon.GetLevel(depth - 1).FindChar('>'), '.')
	return dungeon, explored, Position{playerPos.x, playerPos.y, 0}
}
