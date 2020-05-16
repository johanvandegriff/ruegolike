package main

import "math/rand"

//Generate - generate all the levels in the game
func Generate() ([depth][height][width]int32, [depth][height][width]bool, Position) {
	var levels [depth][height][width]int32
	var explored [depth][height][width]bool

	//simple terrain generation
	var stairX, stairY int
	for z := 0; z < depth; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if rand.Intn(100) < 40 {
					levels[z][y][x] = '#' //wall, 40%
				} else {
					levels[z][y][x] = '.' //empty, 60%
				}
			}
		}
		if z > 0 {
			levels[z][stairY][stairX] = '<'
		}
		if z < depth-1 {
			for ok := true; ok; ok = levels[0][stairY][stairX] != '.' {
				stairX = rand.Intn(width)
				stairY = rand.Intn(height)
			}
			levels[z][stairY][stairX] = '>'
		}
	}
	// level[5][4] = '£'
	// level[5][6] = '#'
	// level[5][6] = '@'

	//start the player on an empty square
	var playerX, playerY int
	//do while
	for ok := true; ok; ok = levels[0][playerY][playerX] != '.' {
		playerX = rand.Intn(width)
		playerY = rand.Intn(height)
	}

	return levels, explored, Position{playerX, playerY, 0}
}
