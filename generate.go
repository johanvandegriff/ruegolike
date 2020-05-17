package main

import (
	"math/rand"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// http://roguebasin.roguelikedevelopment.org/index.php?title=Cellular_Automata_Method_for_Generating_Random_Cave-Like_Levels
func genCaveLevel(level *[height][width]int32) {
	//48 16  40  5 1 4  5 0 3
	fillprob := 40

	//                 r1,r2,reps r1,r2,reps
	gens := [...][3]int{{5, 1, 4}, {5, 0, 3}}

	var grid2 [height][width]int32

	for yi := 0; yi < height; yi++ {
		for xi := 0; xi < width; xi++ {
			if rand.Intn(100) < fillprob {
				level[yi][xi] = '#' //wall, 40%
			} else {
				level[yi][xi] = '.' //empty, 60%
			}
			grid2[yi][xi] = '#'
		}
	}
	//border around the edge
	for yi := 0; yi < height; yi++ {
		level[yi][0] = '#'
		level[yi][width-1] = '#'
	}
	for xi := 0; xi < width; xi++ {
		level[0][xi] = '#'
		level[height-1][xi] = '#'
	}

	for i := 0; i < len(gens); i++ {
		gen := gens[i]
		r1Cutoff := gen[0]
		r2Cutoff := gen[1]
		reps := gen[2]

		for j := 0; j < reps; j++ {

			for yi := 1; yi < height-1; yi++ {
				for xi := 1; xi < width-1; xi++ {
					adjCountR1, adjCountR2 := 0, 0

					for ii := -1; ii <= 1; ii++ {
						for jj := -1; jj <= 1; jj++ {
							if level[yi+ii][xi+jj] != '.' {
								adjCountR1++
							}
						}
					}
					for ii := yi - 2; ii <= yi+2; ii++ {
						for jj := xi - 2; jj <= xi+2; jj++ {
							if abs(ii-yi) == 2 && abs(jj-xi) == 2 {
								continue
							}
							if ii < 0 || jj < 0 || ii >= height || jj >= width {
								continue
							}
							if level[ii][jj] != '.' {
								adjCountR2++
							}
						}
					}
					if adjCountR1 >= r1Cutoff || adjCountR2 <= r2Cutoff {
						grid2[yi][xi] = '#'
					} else {
						grid2[yi][xi] = '.'
					}
				}
			}
			for yi := 1; yi < height-1; yi++ {
				for xi := 1; xi < width-1; xi++ {
					level[yi][xi] = grid2[yi][xi]
				}
			}
		}
	}

}

func floodFill(x, y int, level *[height][width]int32) [height][width]bool {
	var mask [height][width]bool
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			floodFillAux(x+dx, y+dy, level, &mask)
		}
	}
	// floodFillAux(x, y, level, &mask)
	return mask
}

func floodFillAux(x, y int, level *[height][width]int32, mask *[height][width]bool) {
	if isXYInRange(x, y) && level[y][x] == '.' && mask[y][x] == false {
		mask[y][x] = true
		// level[y][x] = '~'
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				floodFillAux(x+dx, y+dy, level, mask)
			}
		}
	}
}

//Generate - generate all the levels in the game
func Generate() ([depth][height][width]int32, [depth][height][width]bool, Position) {
	var levels [depth][height][width]int32
	var explored [depth][height][width]bool

	//simple terrain generation
	var stairX, stairY, playerX, playerY int
	for z := 0; z < depth; z++ {
		genCaveLevel(&levels[z])
		if z > 0 {
			if z == 1 {
				for i := 0; ; i++ {
					if i > 9999 {
						genCaveLevel(&levels[z-1])
						genCaveLevel(&levels[z])
					}
					stairX = rand.Intn(width)
					stairY = rand.Intn(height)
					if levels[z-1][stairY][stairX] == '.' {
						break
					}
				}
				playerX, playerY = stairX, stairY
			}
			mask := floodFill(stairX, stairY, &levels[z-1])
			oldStairX, oldStairY := stairX, stairY
			for i := 0; ; i++ {
				if i > 9999 {
					genCaveLevel(&levels[z])
				}
				stairX = rand.Intn(width)
				stairY = rand.Intn(height)
				if mask[stairY][stairX] && levels[z][stairY][stairX] == '.' &&
					(stairX-oldStairX)*(stairX-oldStairX)+(stairY-oldStairY)*(stairY-oldStairY) >= 8*8 {
					break
				}
			}
			levels[z-1][stairY][stairX] = '>'
			levels[z][stairY][stairX] = '<'
		}
	}

	// https://unicode-search.net/unicode-namesearch.pl?term=BOX%20DRAWINGS
	if debug {
		levels[0][3][5] = '╭' //'┌'
		levels[0][4][5] = '│'
		levels[0][5][5] = '╰' //'└'

		levels[0][3][6] = '─'
		levels[0][5][6] = '─'

		levels[0][3][7] = '╮' //'┐'
		levels[0][4][7] = '│'
		levels[0][5][7] = '╯' //'┘'

		levels[0][10][4] = '£'
		levels[0][10][5] = '#'
		levels[0][10][6] = '@'
	}

	return levels, explored, Position{playerX, playerY, 0}
}
