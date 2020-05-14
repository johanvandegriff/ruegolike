package main

import (
	"fmt"
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

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// invert := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	const width, height = 48, 16
	const offsetX, offsetY = 1, 2

	var level [width][height]int32

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

		//display the map
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				s.SetContent(x+offsetX, y+offsetY, level[x][y], nil, style)
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
