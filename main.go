package main

import (
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	g.movie = movie
}

func (g *game) getClips(clip string) []*clips.Clip {
	clips, err := g.movie.GetClips("game", "fg", clip)
	if err != nil {
		log.Fatal(err)
	}
	return clips
}

func (g *game) setHandlers() {
	button := g.getClips("button")[0]
	button.OnPress(func(id int) {
		g.button = buttonPressed
	})
	button.OnRelease(func(id int) {
		if g.button == buttonPressed {
			g.restart()
		}
	})
	button.OnReleaseOutside(func(id int) {
		if g.button == buttonPressed {
			g.restart()
		}
	})
	icons := g.getClips("icons")
	for y := 0; y < g.c.height; y++ {
		for x := 0; x < g.c.width; x++ {
			px, py := x, y
			icons[y*g.c.width+x].OnPress(func(id int) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.button = buttonEvaluate
				if g.tiles[py][px].marked {
					return
				}
				g.tiles[py][px].pressed = true
				if g.tiles[py][px].open {
					g.forEachNeighbour(px, py, func(x, y int) {
						if !g.tiles[y][x].marked {
							g.tiles[y][x].pressed = true
						}
					})
				}
			})
			icons[y*g.c.width+x].OnLongPress(func(id int) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.onPressTile(px, py, true)
				g.tiles[py][px].pressed = false
			})
			icons[y*g.c.width+x].OnRelease(func(id int) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.button = buttonPlaying
				if g.tiles[py][px].pressed {
					g.onPressTile(px, py, false)
				}
				g.tiles[py][px].pressed = false
			})
			icons[y*g.c.width+x].OnReleaseOutside(func(id int) {
				if g.state == stateWon || g.state == stateLost {
					return
				}
				g.button = buttonPlaying
				g.tiles[py][px].pressed = false
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

func (g *game) onPressTile(x, y int, long bool) {
	if g.state == stateWaiting {
		g.state = statePlaying
		g.time = time.Now().UnixNano()
		g.placeBombs(x, y, g.bombs)
	}
	if !long && g.tiles[y][x].marked {
		return
	}
	if g.tiles[y][x].open {
		if long {
			var marked = 0
			g.forEachNeighbour(x, y, func(x, y int) {
				if g.tiles[y][x].marked {
					marked++
				}
			})
			if g.tiles[y][x].number == marked {
				g.forEachNeighbour(x, y, func(x, y int) {
					if !g.tiles[y][x].marked {
						g.onPressTile(x, y, false)
					}
				})
			}
		}
	} else {
		if long {
			if g.tiles[y][x].marked {
				g.tiles[y][x].marked = false
				g.bombs++
			} else {
				g.tiles[y][x].marked = true
				g.bombs--
			}
		} else {
			g.tiles[y][x].open = true
			g.closed--
			if g.tiles[y][x].bomb {
				g.state = stateLost
				g.button = buttonLost
				return
			}
			if g.tiles[y][x].number == 0 {
				g.forEachNeighbour(x, y, func(x, y int) {
					g.onPressTile(x, y, false)
				})
			}
		}
	}
}

func (g *game) setButton() {
	button := g.getClips("button")[0]
	button.GotoFrame(g.button)
}

func (g *game) setNumbers() {
	bombsDigits := g.getClips("bombs")
	bombs := g.bombs
	if g.state == stateWon {
		bombs = 0
	}
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
			bombsDigits[2-i].GotoFrame(10)
		} else {
			bombsDigits[2-i].GotoFrame(bombs % 10)
		}
		bombs /= 10
	}
	if g.state == statePlaying || g.state == stateWaiting {
		time := int((time.Now().UnixNano() - g.time) / 1000000000)
		if time > 999 {
			time = 999
		}
		timeDigits := g.getClips("time")
		for i := 0; i < 3; i++ {
			timeDigits[2-i].GotoFrame(time % 10)
			time /= 10
		}
	}
}

func (g *game) setTiles() {
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
				icons[y*g.c.width+x].GotoFrame(icon)
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
				icons[y*g.c.width+x].GotoFrame(icon)
			}
		}
	}
}

func (g *game) Update() error {
	if g.movie == nil {
		g.init()
		g.setHandlers()
	}
	if g.state == stateWaiting {
		g.time = time.Now().UnixNano()
	}
	g.setButton()
	g.setNumbers()
	g.setTiles()
	if g.state == statePlaying {
		if g.closed == g.c.bombs {
			g.state = stateWon
			g.button = buttonWon
		}
	}
	return g.movie.Update()
}

func (g *game) Draw(screen *canvas.Image) {
	g.movie.Draw(screen)
}

func newGame(c config) *game {
	g := &game{c: c}
	return g
}

func (g *game) restart() {
	g.button = buttonPlaying
	g.bombs = g.c.bombs
	g.closed = g.c.height * g.c.height
	g.state = stateWaiting
	g.time = time.Now().UnixNano()
	g.tiles = make([][]tile, g.c.height)
	for y := 0; y < g.c.height; y++ {
		g.tiles[y] = make([]tile, g.c.width)
		for x := 0; x < g.c.width; x++ {
			g.tiles[y][x] = tile{}
		}
	}
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

func main() {
	a := app.NewWithID("com.tqdev.fyne-mines")
	a.SetIcon(resourceMinesiconPng)
	w := a.NewWindow("Fyne Mines")
	g := newGame(config{
		scale:   3,
		width:   8,
		height:  8,
		bombs:   10,
		holding: 15,
	})
	g.restart()
	g.init()
	width, height := g.getSize()

	// Main Menu
	//Beginner (8x8, 10 mines), Intermediate (16x16, 40 mines) and Expert (24x24, 99 mines)
	menuItemBeginner := fyne.NewMenuItem("Beginner", func() {})
	menuItemIntermediate := fyne.NewMenuItem("Intermediate", func() {})
	menuItemExpert := fyne.NewMenuItem("Expert", func() {})
	menuFile := fyne.NewMenu("File ", menuItemBeginner, menuItemIntermediate, menuItemExpert)
	menuItemZoom := fyne.NewMenuItem("Zoom ", func() {})
	//menuItemZoom1x := fyne.NewMenuItem("1:1 pixels", func() {})
	menuItemZoom2x := fyne.NewMenuItem("1:2 pixels", func() {})
	menuItemZoom4x := fyne.NewMenuItem("1:4 pixels", func() {})
	menuItemZoom6x := fyne.NewMenuItem("1:6 pixels", func() {})
	//menuItemZoom8x := fyne.NewMenuItem("1:8 pixels", func() {})
	menuItemZoom.ChildMenu = fyne.NewMenu("" /*menuItemZoom1x,*/, menuItemZoom2x, menuItemZoom4x, menuItemZoom6x /*, menuItemZoom8x*/)
	menuView := fyne.NewMenu("View ", menuItemZoom)
	menuItemAbout := fyne.NewMenuItem("About...", func() {
		dialog.ShowInformation("About Fyne Mines v0.0.1", "Author: Maurits van der Schee\n\ngithub.com/mevdschee/fyne-mines", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemAbout)
	mainMenu := fyne.NewMainMenu(menuFile, menuView, menuHelp)
	w.SetMainMenu(mainMenu)
	container := g.movie.GetContainer()
	w.Resize(fyne.NewSize(float32(width), float32(height+26)))
	w.SetContent(container)
	w.SetPadded(false)
	w.SetFixedSize(true)

	//go runGame()
	w.ShowAndRun()
}
