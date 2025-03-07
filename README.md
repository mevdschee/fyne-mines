# Fyne Mines

![screenshot2](screenshot2.png)

Implementation of minesweeper in Go using the [Fyne](https://fyne.io/) GUI library.

### Building

Install Fyne dependencies:

    sudo apt install golang gcc libgl1-mesa-dev xorg-dev

Install go packages:

    go mod download

Run the application:

    go run .

Note that the first build may take several minutes (!).

### Package using fyne-cross

Install fyne-cross using:

    go install github.com/fyne-io/fyne-cross@latest

Now run the package.sh script to build all binaries.

### Graphics and rules

"[Minesweeper X](https://www.curtisbright.com/msx/)" by Curtis Bright is IMHO the best implementation of Minesweeper ever made. He also provided a [skinning system](https://www.curtisbright.com/msx/skins/skinelements.png). For the rules of the game I have been reading the [MinesweeperGame.com](https://minesweepergame.com) website. As a reference I have also looked at the great [Minesweeper Online](https://minesweeperonline.com) implementation in Javascript.

### Links

You can read some background information on creating this game on my blog:

- [Minesweeper written in Go using Fyne](https://tqdev.com/2024-minesweeper-in-go-using-fyne)
- [A 2D puzzle game in Go using Fyne](https://tqdev.com/2024-creating-a-2d-puzzle-game-in-fyne)
