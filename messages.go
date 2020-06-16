package main

import (
	"strings"

	"github.com/gdamore/tcell"
)

const more = "--more--"
const messageAndMoreWidth = 80 //width + 2
const messageWidth = messageAndMoreWidth - len(more) - 1

var x int = 0

//Message - display a message to the player at the top of the screen, managing messages that are too long with --more--
func Message(m string) {
	if len(m) > messageWidth {
		for i := 0; i < len(m); i += messageWidth {
			end := i + messageWidth
			if end > len(m) {
				end = len(m)
			}
			Message(m[i:end])
		}
		return
	}
	if x != 0 { //if not at beginning of line (so some text is showing)
		EmitStr(Screen, x+1, 0, StyleInvert, more)
		x += len(more)
		Screen.Show()
		WaitForSpace()
		ClearMessage()
	}
	EmitStr(Screen, x, 0, StyleDefault, m)
	Screen.Show()
	x += len(m)
}

//ClearMessage - clear the message part of the screen (such as when the player takes an action, or when the next message needs to be displayed)
func ClearMessage() {
	EmitStr(Screen, 0, 0, StyleDefault, strings.Repeat(" ", messageAndMoreWidth))
	Screen.Show()
	x = 0
}

//WaitForSpace - wait for the space key to be pressed
func WaitForSpace() {
	for {
		ev := Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Rune() == ' ' {
				return
			}
		}
	}
}
