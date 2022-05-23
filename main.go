package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const width, height, depth = 48, 16, 32
const offsetX, offsetY = 1, 2 //the offset from the corner of the screen to display the level
const debug = false

//Screen - global reference to the text screen
var Screen tcell.Screen

//StyleDefault - the display colors for most items on screen
var StyleDefault tcell.Style

//StyleNotVisible - the display colors for things the player cannot currently see
var StyleNotVisible tcell.Style

//StyleDebug - the display colors for the unexplored portions, in debugging
var StyleDebug tcell.Style

//StyleInvert - the display colors for inverted text
var StyleInvert tcell.Style

func main() {
	rand.Seed(time.Now().UnixNano())

	// fmt.Println(NewLevel())

	dungeon, explored, playerPos := Generate()
	maxDepth := playerPos.z
	// time.Sleep(1 * time.Second)

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	Screen = s
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

	StyleDefault = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	StyleNotVisible = tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray).Background(tcell.ColorBlack)
	// StyleNotVisible = tcell.StyleDefault.Foreground(tcell.ColorDarkSlateBlue).Background(tcell.ColorBlack)
	// StyleNotVisible = tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)
	StyleDebug = tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Background(tcell.ColorBlack)
	// style3 := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkGray)
	StyleInvert = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

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
			if debug {
				EmitStr(s, 0, 1, StyleDefault, fmt.Sprintf("%c", ev.Rune()))
				EmitStr(s, 2, 1, StyleDefault, fmt.Sprintf("%s                   ", ev.Name()))
			}
			ClearMessage() // clear the messages when a key is pressed

			var deltaX, deltaY int
			if ev.Rune() == '1' || ev.Rune() == ';' {
				deltaX, deltaY = -1, 1
			} else if ev.Rune() == '2' || ev.Name() == "Down" {
				deltaX, deltaY = 0, 1
			} else if ev.Rune() == '3' || ev.Rune() == '\'' {
				deltaX, deltaY = 1, 1
			} else if ev.Rune() == '4' || ev.Name() == "Left" {
				deltaX, deltaY = -1, 0
			} else if ev.Rune() == '5' {
				// deltaX, deltaY = 0, 0 //TODO run
			} else if ev.Rune() == '6' || ev.Name() == "Right" {
				deltaX, deltaY = 1, 0
			} else if ev.Rune() == '7' || ev.Rune() == 'p' {
				deltaX, deltaY = -1, -1
			} else if ev.Rune() == '8' || ev.Name() == "Up" {
				deltaX, deltaY = 0, -1
			} else if ev.Rune() == '9' || ev.Rune() == '[' {
				deltaX, deltaY = 1, -1
			}
			if deltaX != 0 || deltaY != 0 {
				newPlayerX := playerPos.x + deltaX
				newPlayerY := playerPos.y + deltaY

				if newPlayerX >= 0 && newPlayerX < width &&
					newPlayerY >= 0 && newPlayerY < height &&
					!dungeon.GetTile(Position{newPlayerX, newPlayerY, playerPos.z}).IsSolid() {
					playerPos.x = newPlayerX
					playerPos.y = newPlayerY
				} else {
					Message("OOF!! You ran into a wall.")
				}
			} else if playerPos.z < depth-1 && (ev.Rune() == '.' || ev.Rune() == '>') && (debug || dungeon.GetChar(playerPos) == '>') {
				playerPos.z++
				newPos := dungeon.GetLevel(playerPos.z).FindChar('<')
				playerPos.x = newPos.x
				playerPos.y = newPos.y
				Message("You walk down the stairs.")
				Display(s, playerPos, &visible, &explored[playerPos.z], dungeon.GetLevel(playerPos.z))
				if playerPos.z > maxDepth {
					maxDepth = playerPos.z
					Message("This level you have entered looks completely new to you! (temporary message to test the --more-- functionality)")
				}
			} else if playerPos.z > 0 && (ev.Rune() == '.' || ev.Rune() == '<') && (debug || dungeon.GetChar(playerPos) == '<') {
				playerPos.z--
				newPos := dungeon.GetLevel(playerPos.z).FindChar('>')
				playerPos.x = newPos.x
				playerPos.y = newPos.y
				Message("You walk up the stairs.")
			}
			//TODO fast travel command (or click mouse)
		}

		Display(s, playerPos, &visible, &explored[playerPos.z], dungeon.GetLevel(playerPos.z))
	}
}
