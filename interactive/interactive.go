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
	OnMouseDown func()
	OnMouseUp   func()
}

func NewImage(image *canvas.Image, name string) *Image {
	return &Image{image, name, nil, nil}
}

func (i *Image) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}

func (i *Image) MouseDown(ev *desktop.MouseEvent) {
	log.Println("mouse-down: " + i.name)
	log.Printf("mouse-down: %v\n", ev)
	if i.OnMouseDown != nil {
		i.OnMouseDown()
	}
}

func (i *Image) MouseUp(ev *desktop.MouseEvent) {
	log.Println("mouse-up: " + i.name)
	log.Printf("mouse-up: %v\n", ev)
	if i.OnMouseUp != nil {
		i.OnMouseUp()
	}
}
