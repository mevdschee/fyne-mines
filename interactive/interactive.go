package interactive

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Image struct {
	*canvas.Image
	name        string
	isDown      bool
	OnMouseDown func()
	OnMouseUp   func()
}

// ensure Mousable and Hoverable
var _ desktop.Mouseable = (*Image)(nil)
var _ desktop.Hoverable = (*Image)(nil)

func NewImage(image *canvas.Image, name string) *Image {
	return &Image{image, name, false, nil, nil}
}

func (i *Image) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}

func (i *Image) MouseDown(ev *desktop.MouseEvent) {
	log.Println("mouse-down: " + i.name)
	log.Printf("mouse-down: %v\n", ev)
	//if ev.Button = desktop.MouseButtonPrimary {} // left
	//if ev.Button = desktop.MouseButtonSecondary {} // right
	//if ev.Button = desktop.MouseButtonTertiary {} // middle
	if i.OnMouseDown != nil {
		i.OnMouseDown()
	}
	i.isDown = true
}

func (i *Image) MouseUp(ev *desktop.MouseEvent) {
	log.Println("mouse-up: " + i.name)
	log.Printf("mouse-up: %v\n", ev)
	if i.OnMouseUp != nil {
		i.OnMouseUp()
	}
	i.isDown = false
}

func (i *Image) MouseIn(ev *desktop.MouseEvent) {
	//log.Println("Mouse In")
}

func (i *Image) MouseOut() {
	//log.Println("Mouse Out")
	if i.isDown {
		i.MouseUp(&desktop.MouseEvent{})
	}
}

func (i *Image) MouseMoved(ev *desktop.MouseEvent) {
	//log.Println("Mouse Moved")
	//log.Printf("Mouse Moved %v\n", ev.Button)
}
