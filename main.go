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

func oneRay(x0, y0, x1, y1 float64, level [width][height]int32) bool {
	x0 += 0.5
	x1 += 0.5
	y0 += 0.5
	y1 += 0.5

	dx := math.Abs(x1 - x0)
	dy := math.Abs(y1 - y0)

	x := math.Floor(x0)
	y := math.Floor(y0)

	n := 1
	var xInc, yInc int
	var error float64

	if dx == 0 {
		xInc = 0
		error = math.Inf(1)
	} else if x1 > x0 {
		xInc = 1
		n += int(math.Floor(x1) - x)
		error = (math.Floor(x0) + 1 - x0) * dy
	} else {
		xInc = -1
		n += int(x - math.Floor(x1))
		error = (x0 - math.Floor(x0)) * dy
	}

	if dy == 0 {
		yInc = 0
		error -= math.Inf(1)
	} else if y1 > y0 {
		yInc = 1
		n += int(math.Floor(y1) - y)
		error -= (math.Floor(y0) + 1 - y0) * dx
	} else {
		yInc = -1
		n += int(y - math.Floor(y1))
		error -= (y0 - math.Floor(y0)) * dx
	}

	for ; n > 0; n-- {
		if error > 0 {
			y += float64(yInc)
			error -= dx
		} else {
			x += float64(xInc)
			error += dy
		}
		ix, iy := int(x), int(y)
		if ix < 0 || ix >= width || iy < 0 || iy >= height || level[ix][iy] == '#' {
			return false
		}
	}
	return true
}

func raycast(playerX int, playerY int, visible [width][height]bool, explored [width][height]bool, level [width][height]int32, s tcell.Screen, style tcell.Style) ([width][height]bool, [width][height]bool) {

	//calculate visible and explored tiles with raycasting
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// var dx float64
			// if x > playerX {
			// 	dx = -0.5
			// } else {
			// 	dx = 0.5
			// }
			// var dy float64
			// if y > playerY {
			// 	dy = -0.5
			// } else {
			// 	dy = 0.5
			// }
			visible[x][y] = Abs(playerX-x) <= 1 && Abs(playerY-y) <= 1 ||
				oneRay(float64(playerX), float64(playerY), float64(x)-0.5, float64(y)-0.5, level) ||
				oneRay(float64(playerX), float64(playerY), float64(x)+0.5, float64(y)-0.5, level) ||
				oneRay(float64(playerX), float64(playerY), float64(x)-0.5, float64(y)+0.5, level) ||
				oneRay(float64(playerX), float64(playerY), float64(x)+0.5, float64(y)+0.5, level)
			if visible[x][y] {
				explored[x][y] = true
			}
			continue

			if Abs(playerX-x) <= 1 && Abs(playerY-y) <= 1 {
				visible[x][y] = true
				explored[x][y] = true
				continue
			}

			angle := math.Atan2(float64(y-playerY), float64(x-playerX))
			emitStr(s, 0, 1, style, fmt.Sprintf("%f", angle))

			x2, y2 := float64(x), float64(y)
			// x2 -= 0.5 * math.Cos(angle)
			// y2 -= 0.5 * math.Sin(angle)
			for {
				x2 -= 1 * math.Cos(angle)
				y2 -= 1 * math.Sin(angle)
				if math.Abs(x2-float64(playerX)) < 0.6 && math.Abs(y2-float64(playerY)) < 0.6 {
					visible[x][y] = true
					break
				}
				bad := 0
				x2i := int(math.Ceil(x2))
				y2i := int(math.Ceil(y2))
				if x2i < 0 || x2i >= width || y2i < 0 || y2i >= height || level[x2i][y2i] != '.' {
					if x != x2i || y != y2i {
						bad++
					}
				}
				x2i = int(math.Ceil(x2))
				y2i = int(math.Floor(y2))
				if x2i < 0 || x2i >= width || y2i < 0 || y2i >= height || level[x2i][y2i] != '.' {
					if x != x2i || y != y2i {
						bad++
					}
				}
				x2i = int(math.Floor(x2))
				y2i = int(math.Ceil(y2))
				if x2i < 0 || x2i >= width || y2i < 0 || y2i >= height || level[x2i][y2i] != '.' {
					if x != x2i || y != y2i {
						bad++
					}
				}
				x2i = int(math.Floor(x2))
				y2i = int(math.Floor(y2))
				if x2i < 0 || x2i >= width || y2i < 0 || y2i >= height || level[x2i][y2i] != '.' {
					if x != x2i || y != y2i {
						bad++
					}
				}
				if bad > 1 {
					visible[x][y] = false
					break
				}
			}

			// visible[x][y] = angle > 1.5 //x >= playerX-1
			if visible[x][y] {
				explored[x][y] = true
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
					s.SetContent(x+offsetX, y+offsetY, level[x][y], nil, style2) //tmp
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
