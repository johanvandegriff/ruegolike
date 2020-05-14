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

	const width, height = 48, 16
	const offsetX, offsetY = 1, 2

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

		//calculate visible and explored tiles with raycasting
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if math.Abs(float64(x-playerX)) <= 1 && math.Abs(float64(y-playerY)) <= 1 {
					visible[x][y] = true
					explored[x][y] = true
					continue
				}

				angle := math.Atan2(float64(y-playerY), float64(x-playerX))
				emitStr(s, 0, 1, style, fmt.Sprintf("%f", angle))

				x2, y2 := float64(x), float64(y)
				x2 -= 0.5 * math.Cos(angle)
				y2 -= 0.5 * math.Sin(angle)
				for {
					x2 -= 0.1 * math.Cos(angle)
					y2 -= 0.1 * math.Sin(angle)
					if math.Abs(x2-float64(playerX)) < 0.9 && math.Abs(y2-float64(playerY)) < 0.9 {
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
