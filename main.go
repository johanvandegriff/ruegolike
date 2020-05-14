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

	const width, height = 60, 20

	var level [width][height]int32
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(2) == 1 {
				level[x][y] = '#'
			} else {
				level[x][y] = '.'
			}
		}
	}
	// level[5][4] = '£'
	// level[5][6] = '#'
	// level[5][6] = '@'

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
		ev := s.PollEvent()
		switch ev := ev.(type) {
		// case *tcell.EventResize:
		// 	s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			emitStr(s, 0, 0, style, fmt.Sprintf("%s", ev.Name()))
			emitStr(s, 0, 1, style, fmt.Sprintf("%t", ev.Rune() == '4'))
			emitStr(s, 0, 2, style, fmt.Sprintf("%c", ev.Rune()))
			if ev.Rune() == '2' {
				playerY++
			}
			if ev.Rune() == '4' {
				playerX--
			}
			if ev.Rune() == '6' {
				playerX++
			}
			if ev.Rune() == '8' {
				playerY--
			}
		}
		// var x, y int
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if x == playerX && y == playerY {
					// s.SetContent(x, y, '@', nil, tcell.Style.Blink(style, true))
					s.SetContent(x+5, y+3, '@', nil, style)
				} else {
					s.SetContent(x+5, y+3, level[x][y], nil, style)
				}
			}
		}
		s.ShowCursor(playerX+5, playerY+3)
		// s.SetContent(3, 7, tcell.RuneHLine, nil, style)
		// s.SetContent(3, 8, '#', nil, style)
		// drawBox(s, 1, 2, 3, 4, style, 'a')
		s.Show()
		// s.Sync()
		// time.Sleep(2 * time.Second)
	}
}
