package interactive

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Image struct {
	*canvas.Image
	name              string
	OnTapped          func()
	OnTappedSecondary func()
}

func NewImage(image *canvas.Image, name string) *Image {
	return &Image{image, name, nil, nil}
}

func (i *Image) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}

func (i *Image) Tapped(ev *fyne.PointEvent) {
	log.Println("left-click: " + i.name)
	if i.OnTapped != nil {
		i.OnTapped()
	}
}

func (i *Image) TappedSecondary(ev *fyne.PointEvent) {
	log.Println("right-click: " + i.name)
	if i.OnTappedSecondary != nil {
		i.OnTappedSecondary()
	}
}
