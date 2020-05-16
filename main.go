package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const width, height = 48, 16
const offsetX, offsetY = 1, 2

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

	style1 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	// style2 := tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray).Background(tcell.ColorBlack)
	// style1 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	// style2 := tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)
	// style3 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkGray)

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
			emitStr(s, 0, 0, style1, fmt.Sprintf("%c", ev.Rune()))
			emitStr(s, 2, 0, style1, fmt.Sprintf("%s                   ", ev.Name()))
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
				emitStr(s, 15, 0, style1, "    ")
			} else {
				emitStr(s, 15, 0, style1, "oof!")
			}
		}

		display(s, playerX, playerY, &visible, &explored, &level)
	}
}
