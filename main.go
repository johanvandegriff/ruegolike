package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const width, height, depth = 48, 16, 32
const offsetX, offsetY = 1, 2
const debug = false

func main() {
	rand.Seed(time.Now().UnixNano())

	// fmt.Println(NewLevel())

	dungeon, explored, playerPos := Generate()
	// time.Sleep(1 * time.Second)

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

	var visible [height][width]bool
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
			EmitStr(s, 0, 0, style1, fmt.Sprintf("%c", ev.Rune()))
			EmitStr(s, 2, 0, style1, fmt.Sprintf("%s                   ", ev.Name()))
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
			if deltaX != 0 || deltaY != 0 {
				newPlayerX := playerPos.x + deltaX
				newPlayerY := playerPos.y + deltaY

				if newPlayerX >= 0 && newPlayerX < width &&
					newPlayerY >= 0 && newPlayerY < height &&
					!dungeon.GetTile(Position{newPlayerX, newPlayerY, playerPos.z}).isSolid {
					playerPos.x = newPlayerX
					playerPos.y = newPlayerY
					EmitStr(s, 15, 0, style1, "    ")
				} else {
					EmitStr(s, 15, 0, style1, "oof!")
				}
			} else if ev.Rune() == '>' && (debug || dungeon.GetChar(playerPos) == '>') {
				playerPos.z++
			} else if ev.Rune() == '<' && (debug || dungeon.GetChar(playerPos) == '<') {
				playerPos.z--
			}
		}

		Display(s, playerPos, &visible, &explored[playerPos.z], dungeon.GetLevel(playerPos.z))
	}
}
