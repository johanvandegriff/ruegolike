package main

import (
	"github.com/gdamore/tcell"
	"time"
	"math/rand"
	"fmt"
	"os"
)

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

	style := tcell.StyleDefault.
	Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	var level [10][8] int32
	for x:=0; x<10; x++ {
		for y:=0; y<8; y++ {
			level[x][y] = '.'
		}
	}
	level[5][4] = '£';
	level[5][6] = '#';

	s.Clear()
	for x:=0; x<10; x++ {
		for y:=0; y<8; y++ {
			s.SetContent(x, y, level[x][y], nil, style)
		}
	}
	// s.SetContent(3, 7, tcell.RuneHLine, nil, style)
	// s.SetContent(3, 8, '#', nil, style)
	// drawBox(s, 1, 2, 3, 4, style, 'a')
	s.Show()
	time.Sleep(2 * time.Second)
	s.Fini()
}