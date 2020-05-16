package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

const width, height = 48, 16
const offsetX, offsetY = 1, 2

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func Abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
func Abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
func Abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
func Abs8(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

//https://playtechs.blogspot.com/2007/03/raytracing-on-grid.html
func traceLine(x0, y0, x1, y1 float64) [][]int {
	// x0 += 0.5
	// y0 += 0.5
	// x1 += 0.5
	// y1 += 0.5
	// x0 := float64(x0i) + 0.5
	// y0 := float64(y0i) + 0.5
	// x1 := float64(x1i) + 0.5
	// y1 := float64(y1i) + 0.5

	dx := math.Abs(x1 - x0)
	dy := math.Abs(y1 - y0)

	x := int(math.Floor(x0))
	y := int(math.Floor(y0))

	n := 1
	var xInc, yInc int
	var error float64

	if dx == 0 {
		xInc = 0
		error = math.Inf(1)
	} else if x1 > x0 {
		xInc = 1
		n += int(math.Floor(x1)) - x
		error = (math.Floor(x0) + 1 - x0) * dy
	} else {
		xInc = -1
		n += x - int(math.Floor(x1))
		error = (x0 - math.Floor(x0)) * dy
	}

	if dy == 0 {
		yInc = 0
		error -= math.Inf(1)
	} else if y1 > y0 {
		yInc = 1
		n += int(math.Floor(y1)) - y
		error -= (math.Floor(y0) + 1 - y0) * dx
	} else {
		yInc = -1
		n += y - int(math.Floor(y1))
		error -= (y0 - math.Floor(y0)) * dx
	}

	points := make([][]int, n)

	i := 0
	for ; n > 0; n-- {
		// fmt.Println(x, y, error)
		points[i] = make([]int, 2)
		points[i][0] = x
		points[i][1] = y
		i++

		if error > 0 {
			y += yInc
			error -= dx
		} else {
			x += xInc
			error += dy
		}
	}
	return points
}

func traceLineInt(x0, y0, x1, y1 int) [][]int {
	dx := Abs(x1 - x0)
	dy := Abs(y1 - y0)
	x := x0
	y := y0
	n := 1 + dx + dy
	xInc := -1
	if x1 > x0 {
		xInc = 1
	}
	yInc := -1
	if y1 > y0 {
		yInc = 1
	}
	error := dx - dy
	dx *= 2
	dy *= 2

	points := make([][]int, n)
	i := 0
	// changed := false
	for ; n > 0; n-- {
		// fmt.Println(x, y, error)
		points[i] = make([]int, 2)
		points[i][0] = x
		points[i][1] = y
		// if i >= 2 && points[i-2][0] == points[i-1][0] && points[i-1][0] != points[i][0] {
		// 	if !changed {
		// 		points[i][0] = points[i-1][0]
		// 	}
		// 	changed = !changed
		// }
		// if i >= 2 && points[i-2][1] == points[i-1][1] && points[i-1][1] != points[i][1] {
		// 	if !changed {
		// 		points[i][1] = points[i-1][1]
		// 	}
		// 	changed = !changed
		// }
		i++

		if error > 0 {
			x += xInc
			error -= dy
		} else {
			y += yInc
			error += dx
		}
	}
	return points
}

func isXYInRange(x, y int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
}

func findNeighbors(x, y int, level *[width][height]int32) int {
	// var neighbors [4]bool
	var neighbors int = 0
	if isXYInRange(x, y-1) && level[x][y-1] == '#' {
		neighbors += 8
	}
	if isXYInRange(x-1, y) && level[x-1][y] == '#' {
		neighbors += 4
	}
	if isXYInRange(x, y+1) && level[x][y+1] == '#' {
		neighbors += 2
	}
	if isXYInRange(x+1, y) && level[x+1][y] == '#' {
		neighbors++
	}
	// neighbors[1] = isXYInRange(x-1, y) && level[x-1][y] == '#'
	// neighbors[2] = isXYInRange(x, y+1) && level[x-1][y] == '#'
	// neighbors[3] = isXYInRange(x+1, y) && level[x-1][y] == '#'
	return neighbors
}

func isSquareObstructing(x0, y0, x1, y1, x, y int, level *[width][height]int32) bool {
	// return true //TODO tmp
	//pre-checks
	if y1 < y0 {
		x0, y0, x1, y1 = x1, y1, x0, y0 //make sure that y1 >= y0
	}
	if y1 < y || y0 > y ||
		(x0 > x && x1 > x) ||
		(x0 < x && x1 < x) {
		return false
	}
	if x0 == x1 {
		return x0 == x
	}
	if y0 == y1 {
		return y0 == y
	}
	// if m == 1 || m == -1 {
	if (y1 - y0) == (x1 - x0) {
		return x-x0 == y-y0
	}
	if (y1 - y0) == (x0 - x1) {
		return x-x0 == y0-y
	}
	//now, it is guaranteed that x0 != x1, y0 != y1, y1 > y0, m != 1, m != -1
	// return true //TODO tmp

	neighborsToLines := [16][4]uint8{
		//A B  C  D
		{1, 1, 1, 1},
		{2, 3, 1, 1},
		{1, 2, 3, 1},
		{2, 0, 3, 1},
		{1, 1, 2, 3},
		{2, 3, 2, 3},
		{1, 2, 0, 3},
		{2, 0, 0, 3},
		{1, 1, 1, 2},
		{0, 3, 1, 2},
		{1, 2, 3, 2},
		{0, 0, 3, 2},
		{1, 1, 2, 0},
		{0, 3, 2, 0},
		{1, 2, 0, 0},
		{0, 0, 0, 0}}
	lines := neighborsToLines[findNeighbors(x, y, level)]
	xf, yf := float64(x), float64(y)
	m := float64(y1) - yf/float64(x1) - xf
	b := yf - m*xf
	return doLinesIntersect(xf, yf, m, b, lines)

	// for xf := float64(x0) - 0.5; xf <= float64(x0)+0.5; xf += 0.5 {
	// 	for yf := float64(y0) - 0.5; yf <= float64(y0)+0.5; yf += 0.5 {
	// 		m := float64(y1) - yf/float64(x1) - xf
	// 		b := yf - m*xf
	// 		if !doLinesIntersect(xf, yf, m, b, lines) {
	// 			return false
	// 		}
	// 	}
	// }
	// return true
}

func doLinesIntersect(xf, yf, m, b float64, lines [4]uint8) bool {
	if lines[0] == 1 { //A1
		yi := (m*(-yf+.5+xf) + b) / (1 - m)
		if yf-.5 <= yi && yi <= yf {
			return true
		}
	} else if lines[0] == 2 { //A2
		xi := (yf - .5 - b) / m
		if xf <= xi && xi <= xf+.5 {
			return true
		}
	} else if lines[0] == 3 { //A3
		yi := m*(xf+.5) + b
		if yf-.5 <= yi && yi <= yf {
			return true
		}
	}

	if lines[1] == 1 { //B1
		yi := (m*(yf+.5+xf) + b) / (1 - m)
		if yf <= yi && yi <= yf+.5 {
			return true
		}
	} else if lines[1] == 2 { //B2
		yi := m*(xf+.5) + b
		if yf <= yi && yi <= yf+0.5 {
			return true
		}
	} else if lines[1] == 3 { //B3
		xi := (yf + .5 - b) / m
		if xf <= xi && xi <= xf+.5 {
			return true
		}
	}
	if lines[2] == 1 { //C1
		yi := (m*(-yf-.5+xf) + b) / (1 - m)
		if yf <= yi && yi <= yf+.5 {
			return true
		}
	} else if lines[2] == 2 { //C2
		xi := (yf + .5 - b) / m
		if xf-.5 <= xi && xi <= xf {
			return true
		}
	} else if lines[2] == 3 { //C3
		yi := m*(xf-.5) + b
		if yf <= yi && yi <= yf+.5 {
			return true
		}
	}
	if lines[3] == 1 { //D1
		yi := (m*(yf-.5+xf) + b) / (1 - m)
		if yf-.5 <= yi && yi <= yf {
			return true
		}
	} else if lines[3] == 2 { //D2
		yi := m*(xf-.5) + b
		if yf-.5 <= yi && yi <= yf {
			return true
		}
	} else if lines[3] == 3 { //D3
		xi := (yf - .5 - b) / m
		if xf-.5 <= xi && xi <= xf {
			return true
		}
	}

	return false

}

/*func canPlayerSee(playerX, playerY, x, y int, level *[width][height]int32) bool {
	canSee := false
	for startX := float64(x) - 0.5; startX <= float64(x)+0.5; startX += 0.5 {
		for startY := float64(y) - 0.5; startY <= float64(y)+0.5; startY += 0.5 {
			isRayBlocked := false
			points := traceLine(startX, startY, float64(playerX), float64(playerY))

			for index, point := range points {
				px := point[0]
				py := point[1]
				if index != 0 && index != len(points)-1 &&
					isXYInRange(px, py) && level[px][py] == '#' &&
					isSquareObstructing(playerX, playerY, x, y, px, py, level) {
					isRayBlocked = true
					break
				}
			}
			if !isRayBlocked {
				canSee = true
				break
			}
		}
		if canSee {
			break
		}
	}
	return canSee
}*/

func canPlayerSee(playerX, playerY, x, y int, level *[width][height]int32) bool {
	points := traceLineInt(playerX, playerY, x, y)

	for index, point := range points {
		px := point[0]
		py := point[1]
		if index != 0 && index != len(points)-1 && level[px][py] == '#' &&
			isSquareObstructing(playerX, playerY, x, y, px, py, level) {
			return false
		}
	}
	return true
}

func raycast(playerX int, playerY int, visible *[width][height]bool, explored *[width][height]bool, level *[width][height]int32) {
	//calculate visible and explored tiles with raycasting
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if playerX-5 <= x && x <= playerX+5 && playerY-5 <= y && y <= playerY+5 {
				visible[x][y] = canPlayerSee(playerX, playerY, x, y, level)
				if visible[x][y] {
					explored[x][y] = true
				}
			} else {
				visible[x][y] = false
			}
		}
	}
}

func blocksLight(x, y, octant, originX, originY int, level *[width][height]int32) bool {
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
	return nx < 0 || nx >= width || ny < 0 || ny >= height || level[nx][ny] != '.'
}

func setVisible(x, y, octant, originX, originY int, visible *[width][height]bool, explored *[width][height]bool) {
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
	if nx >= 0 && nx < width && ny >= 0 && ny < height {
		visible[nx][ny] = true
		explored[nx][ny] = true
	}
}

type Slope struct {
	y int
	x int
}

func isSlopeGreater(a, b Slope) bool {
	return a.y*b.x > a.x*b.y
}
func isSlopeGreaterOrEqual(a, b Slope) bool {
	return a.y*b.x >= a.x*b.y
}
func isSlopeLess(a, b Slope) bool {
	return a.y*b.x < a.x*b.y
}

// func isSlopeLessOrEqual(a, b Slope) bool {
// 	return a.y*b.x <= a.x*b.y
// }

//calculate visible and explored tiles with shadowcasting
//see http://www.adammil.net/blog/v125_Roguelike_Vision_Algorithms.html#mine
func shadowcast(originX, originY, rangeLimit int, visible *[width][height]bool, explored *[width][height]bool, level *[width][height]int32) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			visible[x][y] = false
		}
	}
	//loop through each octant
	for octant := 0; octant < 8; octant++ {
		shadowcastAux(octant, originX, originY, rangeLimit, 1, Slope{1, 1}, Slope{0, 1}, visible, explored, level)
	}
}

func shadowcastAux(octant, originX, originY, rangeLimit, x int, top, bottom Slope, visible *[width][height]bool, explored *[width][height]bool, level *[width][height]int32) {
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
				if isSlopeGreaterOrEqual(top, Slope{topY*2 + 1, x * 2}) && !blocksLight(x, topY+1, octant, originX, originY, level) {
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
				if isSlopeGreater(top, Slope{topY*2 + 1, ax}) {
					topY++
				}
			}
		}

		var bottomY int
		if bottom.y == 0 { // if bottom == 0/?, then it's hitting the tile at Y=0 dead center. this is special-cased because
			// bottom.Y starts at zero and remains zero as long as it doesn't hit anything, so it's common
			bottomY = 0
		} else { // bottom > 0
			// if bottom.x == 0 { //TODO
			// 	bottom.x = 1
			// }
			bottomY = ((x*2-1)*bottom.y + bottom.x) / (bottom.x * 2) // the tile that the bottom vector enters from the left
			// code below assumes that if a tile is a wall then it's visible, so if the tile contains a wall we have to
			// ensure that the bottom vector actually hits the wall shape. it misses the wall shape if the top-left corner
			// is beveled and bottom >= (bottomY*2+1)/(x*2). finally, the top-left corner is beveled if the tiles to the
			// left and above are clear. we can assume the tile to the left is clear because otherwise the bottom vector
			// would be greater, so we only have to check above
			if isSlopeGreaterOrEqual(bottom, Slope{bottomY*2 + 1, x * 2}) && blocksLight(x, bottomY, octant, originX, originY, level) && !blocksLight(x, bottomY+1, octant, originX, originY, level) {
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
				isVisible := isOpaque || ((y != topY || isSlopeGreater(top, Slope{y*4 - 1, x*4 + 1})) && (y != bottomY || isSlopeLess(bottom, Slope{y*4 + 1, x*4 - 1})))
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
							if isSlopeGreater(top, Slope{ny, nx}) { // we have to maintain the invariant that top > bottom, so the new sector
								// created by adjusting the bottom is only valid if that's the case
								// if we're at the bottom of the column, then just adjust the current sector rather than recursing
								// since there's no chance that this sector can be split in two by a later transition back to clear
								if y == bottomY {
									bottom = Slope{ny, nx}
									break // don't recurse unless necessary
								} else {
									shadowcastAux(octant, originX, originY, rangeLimit, x+1, top, Slope{ny, nx}, visible, explored, level)
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
							if isSlopeGreaterOrEqual(bottom, Slope{ny, nx}) {
								return
							}
							top = Slope{ny, nx}
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

func main() {
	rand.Seed(time.Now().UnixNano())
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	// fmt.Println(s, e)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	e = s.Init()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	// style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	// style2 := tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	style2 := tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)
	style3 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkGray)

	// invert := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	var level [width][height]int32
	var visible, explored [width][height]bool

	//simple terrain generation
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(100) < 40 {
				level[x][y] = '#' //wall, 40%
			} else {
				level[x][y] = '.' //empty, 60%
			}
		}
	}
	// level[5][4] = '£'
	// level[5][6] = '#'
	// level[5][6] = '@'

	//start the player on an empty square
	var playerX, playerY int
	//do while
	for ok := true; ok; ok = level[playerX][playerY] != '.' {
		playerX = rand.Intn(width)
		playerY = rand.Intn(height)
	}

	// fmt.Println(playerX, playerY)
	// time.Sleep(2 * time.Second)

	s.Clear()
	for {
		//player movement
		ev := s.PollEvent()
		switch ev := ev.(type) {
		// case *tcell.EventResize:
		// 	s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			emitStr(s, 0, 0, style, fmt.Sprintf("%c", ev.Rune()))
			emitStr(s, 2, 0, style, fmt.Sprintf("%s                   ", ev.Name()))
			var deltaX, deltaY int
			if ev.Name() == "Left" {
				deltaX, deltaY = -1, 0
			} else if ev.Name() == "Right" {
				deltaX, deltaY = 1, 0
			} else if ev.Name() == "Up" {
				deltaX, deltaY = 0, -1
			} else if ev.Name() == "Down" {
				deltaX, deltaY = 0, 1
			} else if ev.Rune() == '1' {
				deltaX, deltaY = -1, 1
			} else if ev.Rune() == '2' {
				deltaX, deltaY = 0, 1
			} else if ev.Rune() == '3' {
				deltaX, deltaY = 1, 1
			} else if ev.Rune() == '4' {
				deltaX, deltaY = -1, 0
			} else if ev.Rune() == '5' {
				deltaX, deltaY = 0, 0
			} else if ev.Rune() == '6' {
				deltaX, deltaY = 1, 0
			} else if ev.Rune() == '7' {
				deltaX, deltaY = -1, -1
			} else if ev.Rune() == '8' {
				deltaX, deltaY = 0, -1
			} else if ev.Rune() == '9' {
				deltaX, deltaY = 1, -1
			}

			newPlayerX := playerX + deltaX
			newPlayerY := playerY + deltaY

			if newPlayerX >= 0 && newPlayerX < width &&
				newPlayerY >= 0 && newPlayerY < height &&
				level[newPlayerX][newPlayerY] == '.' {
				playerX = newPlayerX
				playerY = newPlayerY
				emitStr(s, 15, 0, style, "    ")
			} else {
				emitStr(s, 15, 0, style, "oof!")
			}
		}

		// raycast(playerX, playerY, &visible, &explored, &level)
		rangeLimit := 5
		shadowcast(playerX, playerY, rangeLimit, &visible, &explored, &level)

		//display the level
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if explored[x][y] {
					if visible[x][y] {
						s.SetContent(x+offsetX, y+offsetY, level[x][y], nil, style)
					} else {
						s.SetContent(x+offsetX, y+offsetY, level[x][y], nil, style2)
					}
				} else {
					// s.SetContent(x+offsetX, y+offsetY, ' ', nil, style)
					s.SetContent(x+offsetX, y+offsetY, level[x][y], nil, style3) //tmp
				}
			}
		}
		// s.SetContent(x, y, '@', nil, tcell.Style.Blink(style, true))
		s.SetContent(playerX+offsetX, playerY+offsetY, '@', nil, style) //display the player
		s.ShowCursor(playerX+offsetX, playerY+offsetY)                  //highlight the player

		// s.SetContent(3, 7, tcell.RuneHLine, nil, style)
		// s.SetContent(3, 8, '#', nil, style)
		// drawBox(s, 1, 2, 3, 4, style, 'a')
		s.Show()
		// s.Sync()
		// time.Sleep(2 * time.Second)
	}
}
