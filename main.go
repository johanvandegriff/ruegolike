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

func findNeighbors(x, y int, level [width][height]int32) int {
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

func isSquareObstructing(x0, y0, x1, y1, x, y int, level [width][height]int32) bool {
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

/*func canPlayerSee(playerX, playerY, x, y int, level [width][height]int32) bool {
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

func canPlayerSee(playerX, playerY, x, y int, level [width][height]int32) bool {
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

func raycast(playerX int, playerY int, visible [width][height]bool, explored [width][height]bool, level [width][height]int32, s tcell.Screen, style tcell.Style) ([width][height]bool, [width][height]bool) {

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
	return visible, explored
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

		visible, explored = raycast(playerX, playerY, visible, explored, level, s, style)

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
