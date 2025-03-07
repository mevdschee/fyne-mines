package main

import (
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/mevdschee/fyne-mines/clips"
	"github.com/mevdschee/fyne-mines/movies"
	"github.com/mevdschee/fyne-mines/sprites"
)

const spriteMapMeta = `
	[{"name":"display","x":28,"y":82,"width":41,"height":25,"count":1},
	{"name":"icons","x":0,"y":0,"width":16,"height":16,"count":17,"grid":9},
	{"name":"digits","x":0,"y":33,"width":11,"height":21,"count":11,"gap":1},
	{"name":"buttons","x":0,"y":55,"width":26,"height":26,"count":5,"gap":1},
	{"name":"controls","x":0,"y":82,"widths":[12,1,12],"heights":[11,1,11],"gap":1},
	{"name":"field","x":0,"y":96,"widths":[12,1,12],"heights":[11,1,11],"gap":1}]`

const movieScenes = `
	[{"name":"game","layers":[{"name":"bg","clips":[
		{"sprite":"controls","x":"0","y":"0","width":"w*16+24","height":"55"},
		{"sprite":"field","x":"0","y":"44","width":"w*16+24","height":"h*16+22"},
		{"sprite":"display","x":"16","y":"15"},
		{"sprite":"display","x":"w*16-33","y":"15"}
	]},{"name":"fg","clips":[
		{"sprite":"digits","name":"bombs","repeat":"3","x":"18+i*13","y":"17"},
		{"sprite":"digits","name":"time","repeat":"3","x":"w*16-31+i*13","y":"17"},
		{"sprite":"buttons","name":"button","x":"(w*16)/2-1","y":"15"},
		{"sprite":"icons","name":"icons","repeat":"w*h","x":"12+(i%w)*16","y":"55+floor(i/w)*16"}
	]}]}]`

type config struct {
	scale   int
	width   int
	height  int
	bombs   int
	holding int
}

type game struct {
	c      config
	movie  *movies.Movie
	button int
	bombs  int
	closed int
	state  int
	time   int64
	tiles  [][]tile
}

type tile struct {
	open    bool
	marked  bool
	bomb    bool
	pressed bool
	number  int
}

const (
	stateWaiting = iota
	statePlaying
	stateWon
	stateLost
)

const (
	buttonPlaying = iota
	buttonEvaluate
	buttonLost
	buttonWon
	buttonPressed
)

const (
	iconEmpty = iota
	iconNumberOne
	iconNumberTwo
	iconNumberThree
	iconNumberFour
	iconNumberFive
	iconNumberSix
	iconNumberSeven
	iconNumberEight
	iconClosed
	iconOpened
	iconBomb
	iconMarked
	iconAnswerNoBomb
	iconAnswerIsBomb
	iconQuestionMark
	iconQuestionPressed
)

var clipCache map[string][]*clips.Clip

func (g *game) getSize() (int, int) {
	return g.c.scale * (g.c.width*16 + 12*2), g.c.scale * (g.c.height*16 + 11*3 + 33)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.getSize()
}

func (g *game) init() {
	spriteMap, err := sprites.NewSpriteMap(resourceWinxpskinPng.Content(), spriteMapMeta)
	if err != nil {
		log.Fatalln(err)
	}
	parameters := map[string]interface{}{
		"w": g.c.width,
		"h": g.c.height,
		"s": g.c.scale,
	}
	movie, err := movies.FromJSON(spriteMap, movieScenes, parameters)
	if err != nil {
		log.Fatalln(err)
	}
	width, height := g.getSize()
	movie.SetSize(width, height)
	g.movie = movie
	clipCache = map[string][]*clips.Clip{}
}

func (g *game) getClips(clip string) []*clips.Clip {
	if clipCache == nil {
		clipCache = map[string][]*clips.Clip{}
	}
	cache, ok := clipCache[clip]
	if ok {
		return cache
	}
	clips, err := g.movie.GetClips("game", "fg", clip)
	if err != nil {
		log.Fatal(err)
	}
	clipCache[clip] = clips
	return clips
}

func (g *game) setHandlers() {
	button := g.getClips("button")[0]
	button.OnPress(func(left, right, middle, alt, control bool) {
		g.button = buttonPressed
		g.updateButton()
	})
	button.OnRelease(func(left, right, middle, alt, control bool) {
		if g.button == buttonPressed {
			g.restart()
		}
	})
	button.OnLeave(func() {
		if g.button == buttonPressed {
			g.restart()
		}
	})
	icons := g.getClips("icons")
	for y := 0; y < g.c.height; y++ {
		for x := 0; x < g.c.width; x++ {
			px, py := x, y
			icons[y*g.c.width+x].OnPress(func(left, right, middle, alt, control bool) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				if !right {
					if !g.tiles[py][px].marked {
						g.button = buttonEvaluate
						g.updateButton()
						g.tiles[py][px].pressed = true
						g.updateTile(px, py)
						if g.tiles[py][px].open {
							g.forEachNeighbour(px, py, func(x, y int) {
								if !g.tiles[y][x].marked {
									g.tiles[y][x].pressed = true
									g.updateTile(x, y)
								}
							})
						}
					}
				}
				if right {
					if !g.tiles[py][px].open {
						if g.tiles[py][px].marked {
							g.tiles[py][px].marked = false
							g.bombs++
						} else {
							g.tiles[py][px].marked = true
							g.bombs--
						}
						g.updateBombDigits()
						g.updateTile(px, py)
					}
				}
			})
			icons[y*g.c.width+x].OnRelease(func(left, right, middle, alt, control bool) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.button = buttonPlaying
				g.updateButton()
				if !right {
					if !g.tiles[py][px].open {
						if g.tiles[py][px].pressed {
							g.tiles[py][px].pressed = false
							g.onPressTile(px, py)
							g.updateAllTiles()
						}
					} else {
						var marks = 0
						g.forEachNeighbour(px, py, func(x, y int) {
							if g.tiles[y][x].marked {
								marks++
							}
						})
						if g.tiles[py][px].number == marks {
							g.forEachNeighbour(px, py, func(x, y int) {
								if !g.tiles[y][x].open && !g.tiles[y][x].marked {
									g.onPressTile(x, y)
								}
							})
							g.updateAllTiles()
						} else {
							g.forEachNeighbour(px, py, func(x, y int) {
								if !g.tiles[y][x].marked {
									g.tiles[y][x].pressed = false
									g.updateTile(x, y)
								}
							})
						}
					}
				}
			})
			icons[y*g.c.width+x].OnEnter(func(left, right, middle, alt, control bool) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				if left {
					g.button = buttonEvaluate
					g.updateButton()
					g.tiles[py][px].pressed = true
					g.updateTile(px, py)
					if g.tiles[py][px].open {
						g.forEachNeighbour(px, py, func(x, y int) {
							if !g.tiles[y][x].marked {
								g.tiles[y][x].pressed = true
								g.updateTile(x, y)
							}
						})
					}
				}
			})
			icons[y*g.c.width+x].OnLeave(func() {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.button = buttonPlaying
				g.updateButton()
				g.tiles[py][px].pressed = false
				if g.tiles[py][px].open {
					g.forEachNeighbour(px, py, func(x, y int) {
						if !g.tiles[y][x].marked {
							g.tiles[y][x].pressed = false
							g.updateTile(x, y)
						}
					})
				}
				g.updateTile(px, py)
			})
		}
	}
}

func (g *game) forEachNeighbour(x, y int, do func(x, y int)) {
	for i := 0; i < 9; i++ {
		dy, dx := i/3-1, i%3-1
		if dy == 0 && dx == 0 {
			continue
		}
		if y+dy < 0 || x+dx < 0 {
			continue
		}
		if y+dy >= g.c.height || x+dx >= g.c.width {
			continue
		}
		do(x+dx, y+dy)
	}
}

func (g *game) onPressTile(x, y int) {
	if g.state == stateWaiting {
		g.state = statePlaying
		g.time = time.Now().UnixNano()
		g.updateTimeDigits()
		g.placeBombs(x, y, g.bombs)
	}
	if !g.tiles[y][x].open && !g.tiles[y][x].marked {
		g.tiles[y][x].open = true
		g.closed--
		if g.tiles[y][x].bomb {
			g.state = stateLost
			g.button = buttonLost
			g.updateButton()
			return
		}
		if g.closed == g.c.bombs {
			g.state = stateWon
			g.button = buttonWon
			g.updateButton()
			g.bombs = 0
			g.updateBombDigits()
			return
		}
		if g.tiles[y][x].number == 0 {
			g.forEachNeighbour(x, y, func(px, py int) {
				g.onPressTile(px, py)
			})
		}
	}
}

func (g *game) updateButton() {
	button := g.getClips("button")[0]
	button.GotoFrame(g.button, true)
}

func (g *game) updateBombDigits() {
	bombsDigits := g.getClips("bombs")
	bombs := g.bombs
	if bombs < -99 {
		bombs = -99
	}
	negative := false
	if bombs < 0 {
		negative = true
		bombs *= -1
	}
	for i := 0; i < 3; i++ {
		if i == 2 && negative {
			bombsDigits[2-i].GotoFrame(10, true)
		} else {
			bombsDigits[2-i].GotoFrame(bombs%10, true)
		}
		bombs /= 10
	}
}

func (g *game) updateTimeDigits() {
	if g.state == statePlaying {
		time := int((time.Now().UnixNano() - g.time) / 1000000000)
		if time > 999 {
			time = 999
		}
		timeDigits := g.getClips("time")
		for i := 0; i < 3; i++ {
			timeDigits[2-i].GotoFrame(time%10, true)
			time /= 10
		}
	}
}

func (g *game) updateAllTiles() {
	icons := g.getClips("icons")
	if g.state == stateWon || g.state == stateLost {
		for y := 0; y < g.c.height; y++ {
			for x := 0; x < g.c.width; x++ {
				icon := iconClosed
				if g.tiles[y][x].open {
					if g.tiles[y][x].bomb {
						icon = iconAnswerIsBomb
					} else {
						icon = g.tiles[y][x].number
					}
				} else {
					if g.tiles[y][x].marked {
						if g.tiles[y][x].bomb {
							icon = iconMarked
						} else {
							icon = iconAnswerNoBomb
						}
					} else {
						if g.tiles[y][x].bomb {
							if g.state == stateWon {
								icon = iconMarked
							} else {
								icon = iconBomb
							}
						}
					}
				}
				icons[y*g.c.width+x].GotoFrame(icon, false)
			}
		}
	} else {
		for y := 0; y < g.c.height; y++ {
			for x := 0; x < g.c.width; x++ {
				icon := iconClosed
				if g.tiles[y][x].open {
					icon = g.tiles[y][x].number
				} else {
					if g.tiles[y][x].marked {
						icon = iconMarked
					} else {
						if g.tiles[y][x].pressed {
							icon = iconEmpty
						}
						if g.state == stateWon {
							icon = iconMarked
						}
					}
				}
				icons[y*g.c.width+x].GotoFrame(icon, false)
			}
		}
	}
	g.movie.GetContainer().Refresh()
}

func (g *game) updateTile(x, y int) {
	icons := g.getClips("icons")
	icon := iconClosed
	if g.tiles[y][x].open {
		icon = g.tiles[y][x].number
	} else {
		if g.tiles[y][x].marked {
			icon = iconMarked
		} else {
			if g.tiles[y][x].pressed {
				icon = iconEmpty
			}
			if g.state == stateWon {
				icon = iconMarked
			}
		}
	}
	icons[y*g.c.width+x].GotoFrame(icon, true)
}

func (g *game) restart() {
	g.state = stateWaiting
	g.button = buttonPlaying
	g.updateButton()
	g.bombs = g.c.bombs
	g.updateBombDigits()
	g.closed = g.c.width * g.c.height
	g.time = time.Now().UnixNano()
	g.updateTimeDigits()
	g.tiles = make([][]tile, g.c.height)
	for y := 0; y < g.c.height; y++ {
		g.tiles[y] = make([]tile, g.c.width)
		for x := 0; x < g.c.width; x++ {
			g.tiles[y][x] = tile{}
		}
	}
	g.updateAllTiles()
}

func (g *game) placeBombs(x, y, bombs int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := bombs
	g.tiles[y][x].bomb = true
	for b > 0 {
		x, y := rng.Intn(g.c.width), rng.Intn(g.c.height)
		if !g.tiles[y][x].bomb {
			g.tiles[y][x].bomb = true
			b--
			g.forEachNeighbour(x, y, func(x, y int) {
				g.tiles[y][x].number++
			})
		}
	}
	g.tiles[y][x].bomb = false
}

func NewGame(config config, window fyne.Window) *game {
	g := &game{c: config}
	g.init()
	g.setHandlers()
	g.restart()
	window.SetContent(g.movie.GetContainer())
	window.Resize(fyne.NewSize(0, 0))
	return g
}

func main() {
	a := app.NewWithID("com.tqdev.fyne-mines")
	a.SetIcon(resourceMinesiconPng)
	w := a.NewWindow("Fyne Mines")
	var g *game
	c := config{
		scale:   2,
		holding: 15,
	}
	menuItemBeginner := fyne.NewMenuItem("Beginner", func() {
		c.width = 9
		c.height = 9
		c.bombs = 10
		g = NewGame(c, w)
	})
	menuItemIntermediate := fyne.NewMenuItem("Intermediate", func() {
		c.width = 16
		c.height = 16
		c.bombs = 40
		g = NewGame(c, w)
	})
	menuItemExpert := fyne.NewMenuItem("Expert", func() {
		c.width = 30
		c.height = 16
		c.bombs = 99
		g = NewGame(c, w)
	})
	menuGame := fyne.NewMenu("Game ", menuItemBeginner, menuItemIntermediate, menuItemExpert)
	menuItemAbout := fyne.NewMenuItem("About...", func() {
		dialog.ShowInformation("About Fyne Mines v1.1.3", "Author: Maurits van der Schee\n\ngithub.com/mevdschee/fyne-mines", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemAbout)
	mainMenu := fyne.NewMainMenu(menuGame, menuHelp)
	w.SetMainMenu(mainMenu)
	w.SetPadded(false)
	menuItemBeginner.Action()
	w.SetContent(g.movie.GetContainer())
	w.SetFixedSize(true)
	go func() {
		for range time.Tick(time.Millisecond * 100) {
			g.updateTimeDigits()
		}
	}()
	w.ShowAndRun()
}
