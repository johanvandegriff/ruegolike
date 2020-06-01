package main

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

//EmitStr - display a string on the screen
func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func isXYInRange(x, y int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
}

func blocksLight(x, y, octant, originX, originY int, level *Level) bool {
	nx, ny := originX, originY
	switch octant {
	case 0:
		nx += x
		ny -= y
	case 1:
		nx += y
		ny -= x
	case 2:
		nx -= y
		ny -= x
	case 3:
		nx -= x
		ny -= y
	case 4:
		nx -= x
		ny += y
	case 5:
		nx -= y
		ny += x
	case 6:
		nx += y
		ny += x
	case 7:
		nx += x
		ny += y
	}
	return !isXYInRange(nx, ny) || level.GetTile(Point{nx, ny}).BlocksLight()
}

//TODO move visible and explored to player/creature
func setVisible(x, y, octant, originX, originY int, visible *[height][width]bool, explored *[height][width]bool) {
	nx, ny := originX, originY
	switch octant {
	case 0:
		nx += x
		ny -= y
	case 1:
		nx += y
		ny -= x
	case 2:
		nx -= y
		ny -= x
	case 3:
		nx -= x
		ny -= y
	case 4:
		nx -= x
		ny += y
	case 5:
		nx -= y
		ny += x
	case 6:
		nx += y
		ny += x
	case 7:
		nx += x
		ny += y
	}
	if isXYInRange(nx, ny) {
		visible[ny][nx] = true
		explored[ny][nx] = true
	}
}

type slope struct {
	y int
	x int
}

func isSlopeGreater(a, b slope) bool {
	return a.y*b.x > a.x*b.y
}
func isSlopeGreaterOrEqual(a, b slope) bool {
	return a.y*b.x >= a.x*b.y
}
func isSlopeLess(a, b slope) bool {
	return a.y*b.x < a.x*b.y
}

// func isSlopeLessOrEqual(a, b slope) bool {
// 	return a.y*b.x <= a.x*b.y
// }

//calculate visible and explored tiles with shadowcasting
//see http://www.adammil.net/blog/v125_Roguelike_Vision_Algorithms.html#mine
func shadowcast(originX, originY, rangeLimit int, visible *[height][width]bool, explored *[height][width]bool, level *Level) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			visible[y][x] = false
		}
	}
	visible[originY][originX] = true
	explored[originY][originX] = true
	//loop through each octant
	for octant := 0; octant < 8; octant++ {
		shadowcastAux(octant, originX, originY, rangeLimit, 1, slope{1, 1}, slope{0, 1}, visible, explored, level)
	}
}

func shadowcastAux(octant, originX, originY, rangeLimit, x int, top, bottom slope, visible *[height][width]bool, explored *[height][width]bool, level *Level) {
	// throughout this function there are references to various parts of tiles. a tile's coordinates refer to its
	// center, and the following diagram shows the parts of the tile and the vectors from the origin that pass through
	// those parts. given a part of a tile with vector u, a vector v passes above it if v > u and below it if v < u
	//    g         center:        y / x
	// a------b   a top left:      (y*2+1) / (x*2-1)   i inner top left:      (y*4+1) / (x*4-1)
	// |  /\  |   b top right:     (y*2+1) / (x*2+1)   j inner top right:     (y*4+1) / (x*4+1)
	// |i/__\j|   c bottom left:   (y*2-1) / (x*2-1)   k inner bottom left:   (y*4-1) / (x*4-1)
	//e|/|  |\|f  d bottom right:  (y*2-1) / (x*2+1)   m inner bottom right:  (y*4-1) / (x*4+1)
	// |\|__|/|   e middle left:   (y*2) / (x*2-1)
	// |k\  /m|   f middle right:  (y*2) / (x*2+1)     a-d are the corners of the tile
	// |  \/  |   g top center:    (y*2+1) / (x*2)     e-h are the corners of the inner (wall) diamond
	// c------d   h bottom center: (y*2-1) / (x*2)     i-m are the corners of the inner square (1/2 tile width)
	//    h
	// for(; x <= (uint)rangeLimit; x++) // (x <= (uint)rangeLimit) == (rangeLimit < 0 || x <= rangeLimit)
	for ; rangeLimit < 0 || x <= rangeLimit; x++ {
		// compute the Y coordinates of the top and bottom of the sector. we maintain that top > bottom
		var topY int
		if top.x == 1 { // if top == ?/1 then it must be 1/1 because 0/1 < top <= 1/1. this is special-cased because top
			topY = x // starts at 1/1 and remains 1/1 as long as it doesn't hit anything, so it's a common case
		} else { // top < 1
			// get the tile that the top vector enters from the left. since our coordinates refer to the center of the
			// tile, this is (x-0.5)*top+0.5, which can be computed as (x-0.5)*top+0.5 = (2(x+0.5)*top+1)/2 =
			// ((2x+1)*top+1)/2. since top == a/b, this is ((2x+1)*a+b)/2b. if it enters a tile at one of the left
			// corners, it will round up, so it'll enter from the bottom-left and never the top-left
			topY = ((x*2-1)*top.y + top.x) / (top.x * 2) // the Y coordinate of the tile entered from the left
			// now it's possible that the vector passes from the left side of the tile up into the tile above before
			// exiting from the right side of this column. so we may need to increment topY
			if blocksLight(x, topY, octant, originX, originY, level) { // if the tile blocks light (i.e. is a wall)...
				// if the tile entered from the left blocks light, whether it passes into the tile above depends on the shape
				// of the wall tile as well as the angle of the vector. if the tile has does not have a beveled top-left
				// corner, then it is blocked. the corner is beveled if the tiles above and to the left are not walls. we can
				// ignore the tile to the left because if it was a wall tile, the top vector must have entered this tile from
				// the bottom-left corner, in which case it can't possibly enter the tile above.
				//
				// otherwise, with a beveled top-left corner, the slope of the vector must be greater than or equal to the
				// slope of the vector to the top center of the tile (x*2, topY*2+1) in order for it to miss the wall and
				// pass into the tile above
				if isSlopeGreaterOrEqual(top, slope{topY*2 + 1, x * 2}) && !blocksLight(x, topY+1, octant, originX, originY, level) {
					topY++
				}
			} else { // the tile doesn't block light
				// since this tile doesn't block light, there's nothing to stop it from passing into the tile above, and it
				// does so if the vector is greater than the vector for the bottom-right corner of the tile above. however,
				// there is one additional consideration. later code in this method assumes that if a tile blocks light then
				// it must be visible, so if the tile above blocks light we have to make sure the light actually impacts the
				// wall shape. now there are three cases: 1) the tile above is clear, in which case the vector must be above
				// the bottom-right corner of the tile above, 2) the tile above blocks light and does not have a beveled
				// bottom-right corner, in which case the vector must be above the bottom-right corner, and 3) the tile above
				// blocks light and does have a beveled bottom-right corner, in which case the vector must be above the
				// bottom center of the tile above (i.e. the corner of the beveled edge).
				//
				// now it's possible to merge 1 and 2 into a single check, and we get the following: if the tile above and to
				// the right is a wall, then the vector must be above the bottom-right corner. otherwise, the vector must be
				// above the bottom center. this works because if the tile above and to the right is a wall, then there are
				// two cases: 1) the tile above is also a wall, in which case we must check against the bottom-right corner,
				// or 2) the tile above is not a wall, in which case the vector passes into it if it's above the bottom-right
				// corner. so either way we use the bottom-right corner in that case. now, if the tile above and to the right
				// is not a wall, then we again have two cases: 1) the tile above is a wall with a beveled edge, in which
				// case we must check against the bottom center, or 2) the tile above is not a wall, in which case it will
				// only be visible if light passes through the inner square, and the inner square is guaranteed to be no
				// larger than a wall diamond, so if it wouldn't pass through a wall diamond then it can't be visible, so
				// there's no point in incrementing topY even if light passes through the corner of the tile above. so we
				// might as well use the bottom center for both cases.
				ax := x * 2                                                    // center
				if blocksLight(x+1, topY+1, octant, originX, originY, level) { // use bottom-right if the tile above and right is a wall
					ax++
				}
				if isSlopeGreater(top, slope{topY*2 + 1, ax}) {
					topY++
				}
			}
		}

		var bottomY int
		if bottom.y == 0 { // if bottom == 0/?, then it's hitting the tile at Y=0 dead center. this is special-cased because
			// bottom.Y starts at zero and remains zero as long as it doesn't hit anything, so it's common
			bottomY = 0
		} else { // bottom > 0
			bottomY = ((x*2-1)*bottom.y + bottom.x) / (bottom.x * 2) // the tile that the bottom vector enters from the left
			// code below assumes that if a tile is a wall then it's visible, so if the tile contains a wall we have to
			// ensure that the bottom vector actually hits the wall shape. it misses the wall shape if the top-left corner
			// is beveled and bottom >= (bottomY*2+1)/(x*2). finally, the top-left corner is beveled if the tiles to the
			// left and above are clear. we can assume the tile to the left is clear because otherwise the bottom vector
			// would be greater, so we only have to check above
			if isSlopeGreaterOrEqual(bottom, slope{bottomY*2 + 1, x * 2}) && blocksLight(x, bottomY, octant, originX, originY, level) && !blocksLight(x, bottomY+1, octant, originX, originY, level) {
				bottomY++
			}
		}

		// go through the tiles in the column now that we know which ones could possibly be visible
		wasOpaque := -1 // 0:false, 1:true, -1:not applicable
		// for(uint y = topY; (int)y >= (int)bottomY; y--) // use a signed comparison because y can wrap around when decremented
		for y := topY; y >= bottomY; y-- {
			if rangeLimit < 0 || math.Sqrt(float64(x*x+y*y)) <= float64(rangeLimit) { // skip the tile if it's out of visual range
				isOpaque := blocksLight(x, y, octant, originX, originY, level)
				// every tile where topY > y > bottomY is guaranteed to be visible. also, the code that initializes topY and
				// bottomY guarantees that if the tile is opaque then it's visible. so we only have to do extra work for the
				// case where the tile is clear and y == topY or y == bottomY. if y == topY then we have to make sure that
				// the top vector is above the bottom-right corner of the inner square. if y == bottomY then we have to make
				// sure that the bottom vector is below the top-left corner of the inner square
				isVisible := isOpaque || ((y != topY || isSlopeGreater(top, slope{y*4 - 1, x*4 + 1})) && (y != bottomY || isSlopeLess(bottom, slope{y*4 + 1, x*4 - 1})))
				// NOTE: if you want the algorithm to be either fully or mostly symmetrical, replace the line above with the
				// following line (and uncomment the Slope.LessOrEqual method). the line ensures that a clear tile is visible
				// only if there's an unobstructed line to its center. if you want it to be fully symmetrical, also remove
				// the "isOpaque ||" part and see NOTE comments further down
				// bool isVisible = isOpaque || ((y != topY || top.GreaterOrEqual(y, x)) && (y != bottomY || bottom.LessOrEqual(y, x)));
				if isVisible {
					setVisible(x, y, octant, originX, originY, visible, explored)
				}

				// if we found a transition from clear to opaque or vice versa, adjust the top and bottom vectors
				if x != rangeLimit { // but don't bother adjusting them if this is the last column anyway
					if isOpaque {
						if wasOpaque == 0 { // if we found a transition from clear to opaque, this sector is done in this column,
							// so adjust the bottom vector upward and continue processing it in the next column
							// if the opaque tile has a beveled top-left corner, move the bottom vector up to the top center.
							// otherwise, move it up to the top left. the corner is beveled if the tiles above and to the left are
							// clear. we can assume the tile to the left is clear because otherwise the vector would be higher, so
							// we only have to check the tile above
							nx, ny := x*2, y*2+1 // top center by default
							// NOTE: if you're using full symmetry and want more expansive walls (recommended), comment out the next line
							if blocksLight(x, y+1, octant, originX, originY, level) {
								nx--
							} // top left if the corner is not beveled
							if isSlopeGreater(top, slope{ny, nx}) { // we have to maintain the invariant that top > bottom, so the new sector
								// created by adjusting the bottom is only valid if that's the case
								// if we're at the bottom of the column, then just adjust the current sector rather than recursing
								// since there's no chance that this sector can be split in two by a later transition back to clear
								if y == bottomY {
									bottom = slope{ny, nx}
									break // don't recurse unless necessary
								} else {
									shadowcastAux(octant, originX, originY, rangeLimit, x+1, top, slope{ny, nx}, visible, explored, level)
								}
							} else { // the new bottom is greater than or equal to the top, so the new sector is empty and we'll ignore
								// it. if we're at the bottom of the column, we'd normally adjust the current sector rather than
								if y == bottomY { // recursing, so that invalidates the current sector and we're done
									return
								}
							}
						}
						wasOpaque = 1
					} else {
						if wasOpaque > 0 { // if we found a transition from opaque to clear, adjust the top vector downwards
							// if the opaque tile has a beveled bottom-right corner, move the top vector down to the bottom center.
							// otherwise, move it down to the bottom right. the corner is beveled if the tiles below and to the right
							// are clear. we know the tile below is clear because that's the current tile, so just check to the right
							nx, ny := x*2, y*2+1 // the bottom of the opaque tile (oy*2-1) equals the top of this tile (y*2+1)
							// NOTE: if you're using full symmetry and want more expansive walls (recommended), comment out the next line
							if blocksLight(x+1, y+1, octant, originX, originY, level) {
								nx++
							} // check the right of the opaque tile (y+1), not this one
							// we have to maintain the invariant that top > bottom. if not, the sector is empty and we're done
							if isSlopeGreaterOrEqual(bottom, slope{ny, nx}) {
								return
							}
							top = slope{ny, nx}
						}
						wasOpaque = 0
					}
				}
			}
		}

		// if the column didn't end in a clear tile, then there's no reason to continue processing the current sector
		// because that means either 1) wasOpaque == -1, implying that the sector is empty or at its range limit, or 2)
		// wasOpaque == 1, implying that we found a transition from clear to opaque and we recursed and we never found
		// a transition back to clear, so there's nothing else for us to do that the recursive method hasn't already. (if
		// we didn't recurse (because y == bottomY), it would have executed a break, leaving wasOpaque equal to 0.)
		if wasOpaque != 0 {
			break
		}
	}

}

//Display - display the current level on the ascii screen
func Display(s tcell.Screen, playerPos Position, visible *[height][width]bool, explored1 *[height][width]bool, level *Level) {
	playerX := playerPos.x
	playerY := playerPos.y

	style1 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	style2 := tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray).Background(tcell.ColorBlack)

	EmitStr(s, 0, 19, style1, fmt.Sprintf("coords:%02d,%02d,%02d", playerPos.x, playerPos.y, playerPos.z))

	rangeLimit := -1
	shadowcast(playerX, playerY, rangeLimit, visible, explored1, level)

	//display the level
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if explored1[y][x] {
				char := level.GetChar(Point{x, y})
				if level.GetTile(Point{x, y}).IsRoom() && char != '*' {
					var isFloor [3][3]bool
					numFloors := 0
					for x2 := -1; x2 <= 1; x2++ {
						for y2 := -1; y2 <= 1; y2++ {
							isFloor[x2+1][y2+1] = isXYInRange(x+x2, y+y2) && level.GetTile(Point{x + x2, y + y2}).IsRoomFloor() && explored1[y+y2][x+x2]
							if isFloor[x2+1][y2+1] {
								numFloors++
							}
						}
					}
					if numFloors == 0 {
						char = ' '
					} else {
						// var isRoom [3][3]bool
						// for x2 := -1; x2 <= 1; x2++ {
						// 	for y2 := -1; y2 <= 1; y2++ {
						// 		isRoom[x2+1][y2+1] = isXYInRange(x+x2, y+y2) && level.GetTile(Point{x + x2, y + y2}).IsRoom() && explored1[y+y2][x+x2]
						// 	}
						// }
						isRoomAbove := isXYInRange(x, y-1) && level.GetTile(Point{x, y - 1}).IsRoom() && explored1[y-1][x]
						isRoomBelow := isXYInRange(x, y+1) && level.GetTile(Point{x, y + 1}).IsRoom() && explored1[y+1][x]
						isRoomLeft := isXYInRange(x-1, y) && level.GetTile(Point{x - 1, y}).IsRoom() && explored1[y][x-1]
						isRoomRight := isXYInRange(x+1, y) && level.GetTile(Point{x + 1, y}).IsRoom() && explored1[y][x+1]

						//corner cases :)
						if char == '┌' {
							if isRoomBelow && isRoomRight {
								// char = '┌'
							} else if !isRoomBelow && isRoomRight {
								char = '╶'
							} else if isRoomBelow && !isRoomRight {
								char = '╷'
							} else if !isRoomBelow && !isRoomRight {
								char = ' '
							}
						} else if char == '└' {
							if isRoomAbove && isRoomRight {
								// char = '└'
							} else if !isRoomAbove && isRoomRight {
								char = '╶'
							} else if isRoomAbove && !isRoomRight {
								char = '╵'
							} else if !isRoomAbove && !isRoomRight {
								char = ' '
							}
						} else if char == '┐' {
							if isRoomBelow && isRoomLeft {
								// char = '┐'
							} else if !isRoomBelow && isRoomLeft {
								char = '╴'
							} else if isRoomBelow && !isRoomLeft {
								char = '╷'
							} else if !isRoomBelow && !isRoomLeft {
								char = ' '
							}
						} else if char == '┘' {
							if isRoomAbove && isRoomLeft {
								// char = '┘'
							} else if !isRoomAbove && isRoomLeft {
								char = '╴'
							} else if isRoomAbove && !isRoomLeft {
								char = '╵'
							} else if !isRoomAbove && !isRoomLeft {
								char = ' '
							}
						} else if char == '─' {
							if isRoomLeft && isRoomRight {
								// char = '─'
							} else if !isRoomLeft && isRoomRight {
								char = '╶'
							} else if isRoomLeft && !isRoomRight {
								char = '╴'
							} else if !isRoomLeft && !isRoomRight {
								char = ' '
							}
						} else if char == '│' {
							if isRoomAbove && isRoomBelow {
								// char = '│'
							} else if !isRoomAbove && isRoomBelow {
								char = '╷'
							} else if isRoomAbove && !isRoomBelow {
								char = '╵'
							} else if !isRoomAbove && !isRoomBelow {
								char = ' '
							}
						} else if char == '├' {
							if isRoomAbove && isRoomBelow && isRoomRight {
								// char = '├'
							} else if !isRoomAbove && isRoomBelow && isRoomRight {
								char = '┌'
							} else if isRoomAbove && !isRoomBelow && isRoomRight {
								char = '└'
							} else if !isRoomAbove && !isRoomBelow && isRoomRight {
								char = '╶'
							} else if isRoomAbove && isRoomBelow && !isRoomRight {
								char = '│'
							} else if !isRoomAbove && isRoomBelow && !isRoomRight {
								char = '╷'
							} else if isRoomAbove && !isRoomBelow && !isRoomRight {
								char = '╵'
							} else if !isRoomAbove && !isRoomBelow && !isRoomRight {
								char = ' '
							}
						} else if char == '┤' {
							if isRoomAbove && isRoomBelow && isRoomLeft {
								// char = '┤'
							} else if !isRoomAbove && isRoomBelow && isRoomLeft {
								char = '┐'
							} else if isRoomAbove && !isRoomBelow && isRoomLeft {
								char = '┘'
							} else if !isRoomAbove && !isRoomBelow && isRoomLeft {
								char = '╴'
							} else if isRoomAbove && isRoomBelow && !isRoomLeft {
								char = '│'
							} else if !isRoomAbove && isRoomBelow && !isRoomLeft {
								char = '╷'
							} else if isRoomAbove && !isRoomBelow && !isRoomLeft {
								char = '╵'
							} else if !isRoomAbove && !isRoomBelow && !isRoomLeft {
								char = ' '
							}
						} else if char == '┬' {
							if isRoomLeft && isRoomRight && isRoomBelow {
								// char = '┬'
							} else if !isRoomLeft && isRoomRight && isRoomBelow {
								char = '┌'
							} else if isRoomLeft && !isRoomRight && isRoomBelow {
								char = '┐'
							} else if !isRoomLeft && !isRoomRight && isRoomBelow {
								char = '╷'
							} else if isRoomLeft && isRoomRight && !isRoomBelow {
								char = '─'
							} else if !isRoomLeft && isRoomRight && !isRoomBelow {
								char = '╶'
							} else if isRoomLeft && !isRoomRight && !isRoomBelow {
								char = '╴'
							} else if !isRoomLeft && !isRoomRight && !isRoomBelow {
								char = ' '
							}
						} else if char == '┴' {
							if isRoomLeft && isRoomRight && isRoomAbove {
								// char = '┴'
							} else if !isRoomLeft && isRoomRight && isRoomAbove {
								char = '└'
							} else if isRoomLeft && !isRoomRight && isRoomAbove {
								char = '┘'
							} else if !isRoomLeft && !isRoomRight && isRoomAbove {
								char = '╵'
							} else if isRoomLeft && isRoomRight && !isRoomAbove {
								char = '─'
							} else if !isRoomLeft && isRoomRight && !isRoomAbove {
								char = '╶'
							} else if isRoomLeft && !isRoomRight && !isRoomAbove {
								char = '╴'
							} else if !isRoomLeft && !isRoomRight && !isRoomAbove {
								char = ' '
							}
						} else if char == '┼' {
							if isRoomLeft && isRoomRight && isRoomAbove && isRoomBelow {
								// char = '┼'
							} else if !isRoomLeft && isRoomRight && isRoomAbove && isRoomBelow {
								char = '├'
							} else if isRoomLeft && !isRoomRight && isRoomAbove && isRoomBelow {
								char = '┤'
							} else if !isRoomLeft && !isRoomRight && isRoomAbove && isRoomBelow {
								char = '│'
							} else if isRoomLeft && isRoomRight && !isRoomAbove && isRoomBelow {
								char = '┬'
							} else if !isRoomLeft && isRoomRight && !isRoomAbove && isRoomBelow {
								char = '┌'
							} else if isRoomLeft && !isRoomRight && !isRoomAbove && isRoomBelow {
								char = '┐'
							} else if !isRoomLeft && !isRoomRight && !isRoomAbove && isRoomBelow {
								char = '╷'
							} else if isRoomLeft && isRoomRight && isRoomAbove && !isRoomBelow {
								char = '┴'
							} else if !isRoomLeft && isRoomRight && isRoomAbove && !isRoomBelow {
								char = '└'
							} else if isRoomLeft && !isRoomRight && isRoomAbove && !isRoomBelow {
								char = '┘'
							} else if !isRoomLeft && !isRoomRight && isRoomAbove && !isRoomBelow {
								char = '╵'
							} else if isRoomLeft && isRoomRight && !isRoomAbove && !isRoomBelow {
								char = '─'
							} else if !isRoomLeft && isRoomRight && !isRoomAbove && !isRoomBelow {
								char = '╶'
							} else if isRoomLeft && !isRoomRight && !isRoomAbove && !isRoomBelow {
								char = '╴'
							} else if !isRoomLeft && !isRoomRight && !isRoomAbove && !isRoomBelow {
								char = ' '
							}
						}

						isFloorTop := isFloor[1][0]
						isFloorBottom := isFloor[1][2]
						isFloorLeft := isFloor[0][1]
						isFloorRight := isFloor[2][1]

						isFloorTopLeft := isFloor[0][0]
						isFloorTopRight := isFloor[2][0]
						isFloorBottomLeft := isFloor[0][2]
						isFloorBottomRight := isFloor[2][2]

						isFloorTopAny := isFloorTop || isFloorTopLeft || isFloorTopRight
						isFloorBottomAny := isFloorBottom || isFloorBottomLeft || isFloorBottomRight
						isFloorLeftAny := isFloorLeft || isFloorTopLeft || isFloorBottomLeft
						isFloorRightAny := isFloorRight || isFloorTopRight || isFloorBottomRight

						if char == '┌' {
							if !isFloorBottomRight {
								char = ' '
							}
						} else if char == '└' {
							if !isFloorTopRight {
								char = ' '
							}
						} else if char == '┐' {
							if !isFloorBottomLeft {
								char = ' '
							}
						} else if char == '┘' {
							if !isFloorTopLeft {
								char = ' '
							}
						} else if char == '─' {
							if !isFloorTopAny && !isFloorBottomAny {
								char = ' '
							}
						} else if char == '│' {
							if !isFloorLeftAny && !isFloorRightAny {
								char = ' '
							}
						} else if char == '├' {
							if !isFloorLeftAny && !isFloorTopRight && !isFloorBottomRight {
								char = ' '
							} else if !isFloorLeftAny && !isFloorTopRight && isFloorBottomRight {
								char = '┌'
							} else if !isFloorLeftAny && isFloorTopRight && !isFloorBottomRight {
								char = '└'
							} else if !isFloorLeftAny && isFloorTopRight && isFloorBottomRight {
								// char = '├'
							} else if isFloorLeftAny && !isFloorTopRight && !isFloorBottomRight {
								char = '│'
							} else if isFloorLeftAny && !isFloorTopRight && isFloorBottomRight {
								// char = '├'
							} else if isFloorLeftAny && isFloorTopRight && !isFloorBottomRight {
								// char = '├'
							} else if isFloorLeftAny && isFloorTopRight && isFloorBottomRight {
								// char = '├'
							}
						} else if char == '┤' {
							if !isFloorRightAny && !isFloorTopLeft && !isFloorBottomLeft {
								char = ' '
							} else if !isFloorRightAny && !isFloorTopLeft && isFloorBottomLeft {
								char = '┐'
							} else if !isFloorRightAny && isFloorTopLeft && !isFloorBottomLeft {
								char = '┘'
							} else if !isFloorRightAny && isFloorTopLeft && isFloorBottomLeft {
								// char = '┤'
							} else if isFloorRightAny && !isFloorTopLeft && !isFloorBottomLeft {
								char = '│'
							} else if isFloorRightAny && !isFloorTopLeft && isFloorBottomLeft {
								// char = '┤'
							} else if isFloorRightAny && isFloorTopLeft && !isFloorBottomLeft {
								// char = '┤'
							} else if isFloorRightAny && isFloorTopLeft && isFloorBottomLeft {
								// char = '┤'
							}
						} else if char == '┬' {
							if !isFloorTopAny && !isFloorBottomLeft && !isFloorBottomRight {
								char = ' '
							} else if !isFloorTopAny && !isFloorBottomLeft && isFloorBottomRight {
								char = '┌'
							} else if !isFloorTopAny && isFloorBottomLeft && !isFloorBottomRight {
								char = '┐'
							} else if !isFloorTopAny && isFloorBottomLeft && isFloorBottomRight {
								// char = '┬'
							} else if isFloorTopAny && !isFloorBottomLeft && !isFloorBottomRight {
								char = '─'
							} else if isFloorTopAny && !isFloorBottomLeft && isFloorBottomRight {
								// char = '┬'
							} else if isFloorTopAny && isFloorBottomLeft && !isFloorBottomRight {
								// char = '┬'
							} else if isFloorTopAny && isFloorBottomLeft && isFloorBottomRight {
								// char = '┬'
							}
						} else if char == '┴' {
							if !isFloorBottomAny && !isFloorTopLeft && !isFloorTopRight {
								char = ' '
							} else if !isFloorBottomAny && !isFloorTopLeft && isFloorTopRight {
								char = '└'
							} else if !isFloorBottomAny && isFloorTopLeft && !isFloorTopRight {
								char = '┘'
							} else if !isFloorBottomAny && isFloorTopLeft && isFloorTopRight {
								// char = '┴'
							} else if isFloorBottomAny && !isFloorTopLeft && !isFloorTopRight {
								char = '─'
							} else if isFloorBottomAny && !isFloorTopLeft && isFloorTopRight {
								// char = '┴'
							} else if isFloorBottomAny && isFloorTopLeft && !isFloorTopRight {
								// char = '┴'
							} else if isFloorBottomAny && isFloorTopLeft && isFloorTopRight {
								// char = '┴'
							}
						} else if char == '┼' {
							if !isFloorBottomLeft && !isFloorBottomRight && !isFloorTopLeft && !isFloorTopRight {
								char = ' '
							} else if !isFloorBottomLeft && !isFloorBottomRight && !isFloorTopLeft && isFloorTopRight {
								char = '└'
							} else if !isFloorBottomLeft && !isFloorBottomRight && isFloorTopLeft && !isFloorTopRight {
								char = '┘'
							} else if !isFloorBottomLeft && !isFloorBottomRight && isFloorTopLeft && isFloorTopRight {
								char = '┴'
							} else if !isFloorBottomLeft && isFloorBottomRight && !isFloorTopLeft && !isFloorTopRight {
								char = '┌'
							} else if !isFloorBottomLeft && isFloorBottomRight && !isFloorTopLeft && isFloorTopRight {
								char = '├'
							} else if !isFloorBottomLeft && isFloorBottomRight && isFloorTopLeft && !isFloorTopRight {
								// char = '┼'
							} else if !isFloorBottomLeft && isFloorBottomRight && isFloorTopLeft && isFloorTopRight {
								// char = '┼'
							} else if isFloorBottomLeft && !isFloorBottomRight && !isFloorTopLeft && !isFloorTopRight {
								char = '┐'
							} else if isFloorBottomLeft && !isFloorBottomRight && !isFloorTopLeft && isFloorTopRight {
								// char = '┼'
							} else if isFloorBottomLeft && !isFloorBottomRight && isFloorTopLeft && !isFloorTopRight {
								char = '┤'
							} else if isFloorBottomLeft && !isFloorBottomRight && isFloorTopLeft && isFloorTopRight {
								// char = '┼'
							} else if isFloorBottomLeft && isFloorBottomRight && !isFloorTopLeft && !isFloorTopRight {
								char = '┬'
							} else if isFloorBottomLeft && isFloorBottomRight && !isFloorTopLeft && isFloorTopRight {
								// char = '┼'
							} else if isFloorBottomLeft && isFloorBottomRight && isFloorTopLeft && !isFloorTopRight {
								// char = '┼'
							} else if isFloorBottomLeft && isFloorBottomRight && isFloorTopLeft && isFloorTopRight {
								// char = '┼'
							}
						}
					}
				}
				if visible[y][x] {
					s.SetContent(x+offsetX, y+offsetY, char, nil, style1)
				} else {
					s.SetContent(x+offsetX, y+offsetY, char, nil, style2)
				}
			} else {
				if debug {
					s.SetContent(x+offsetX, y+offsetY, level.GetChar(Point{x, y}), nil, tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Background(tcell.ColorBlack))
				} else {
					s.SetContent(x+offsetX, y+offsetY, ' ', nil, style2)
				}
			}
		}
	}
	// s.SetContent(x, y, '@', nil, tcell.Style.Blink(style, true))
	s.SetContent(playerX+offsetX, playerY+offsetY, '@', nil, style1) //display the player
	s.ShowCursor(playerX+offsetX, playerY+offsetY)                   //highlight the player

	// s.SetContent(3, 7, tcell.RuneHLine, nil, style)
	// s.SetContent(3, 8, ' ', nil, style)
	// drawBox(s, 1, 2, 3, 4, style, 'a')
	s.Show()
	// s.Sync()
	// time.Sleep(2 * time.Second)
}
