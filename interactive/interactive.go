package interactive

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Image struct {
	*canvas.Image
	onMouseDown  func(ev *desktop.MouseEvent)
	onMouseUp    func(ev *desktop.MouseEvent)
	onMouseIn    func(ev *desktop.MouseEvent)
	onMouseOut   func()
	onMouseMoved func(ev *desktop.MouseEvent)
}

// ensure Mousable and Hoverable
var _ desktop.Mouseable = (*Image)(nil)
var _ desktop.Hoverable = (*Image)(nil)

func NewImage(image *canvas.Image) *Image {
	return &Image{Image: image}
}

func (i *Image) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}

// OnMouseDown sets the mouse down handler
func (i *Image) OnMouseDown(handler func(ev *desktop.MouseEvent)) {
	i.onMouseDown = handler
}

func (i *Image) MouseDown(ev *desktop.MouseEvent) {
	if i.onMouseDown != nil {
		i.onMouseDown(ev)
	}
}

// OnMouseUp sets the mouse up handler
func (i *Image) OnMouseUp(handler func(ev *desktop.MouseEvent)) {
	i.onMouseUp = handler
}

func (i *Image) MouseUp(ev *desktop.MouseEvent) {
	if i.onMouseUp != nil {
		i.onMouseUp(ev)
	}
}

// OnMouseIn sets the mouse in handler
func (i *Image) OnMouseIn(handler func(ev *desktop.MouseEvent)) {
	i.onMouseIn = handler
}

func (i *Image) MouseIn(ev *desktop.MouseEvent) {
	if i.onMouseIn != nil {
		i.onMouseIn(ev)
	}
}

// OnMouseOut sets the mouse out handler
func (i *Image) OnMouseOut(handler func()) {
	i.onMouseOut = handler
}

func (i *Image) MouseOut() {
	if i.onMouseOut != nil {
		i.onMouseOut()
	}
}

// OnMouseMoved sets the mouse moved handler
func (i *Image) OnMouseMoved(handler func(ev *desktop.MouseEvent)) {
	i.onMouseMoved = handler
}

func (i *Image) MouseMoved(ev *desktop.MouseEvent) {
	if i.onMouseMoved != nil {
		i.onMouseMoved(ev)
	}
}
