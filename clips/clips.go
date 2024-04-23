package clips

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/mevdschee/fyne-mines/interactive"
	"github.com/mevdschee/fyne-mines/sprites"
	"golang.org/x/image/draw"
)

// Clip is a set of frames
type Clip struct {
	container     *fyne.Container
	name          string
	x, y          int
	width, height int
	scale         int
	overlay       *interactive.Image
	frame         int
	frames        []*canvas.Image
	onPress       func(left, right, middle, alt, control bool)
	onRelease     func(left, right, middle, alt, control bool)
	onEnter       func(left, right, middle, alt, control bool)
	onLeave       func()
	onOver        func(left, right, middle, alt, control bool)
}

// ClipJSON is a clip in JSON
type ClipJSON struct {
	Name          string
	Sprite        string
	Repeat        string
	X, Y          string
	Width, Height string
}

// GetName gets the name of the clip
func (c *Clip) GetName() string {
	return c.name
}

// GetContainer gets the container from the clip
func (c *Clip) GetContainer() *fyne.Container {
	return c.container
}

// GetPosition gets the position of the clip
func (c *Clip) GetPosition() fyne.Position {
	return fyne.Position{X: float32(c.x * c.scale), Y: float32(c.y * c.scale)}
}

// GetSize gets the size of the clip
func (c *Clip) GetSize() fyne.Size {
	return fyne.Size{Width: float32(c.width * c.scale), Height: float32(c.height * c.scale)}
}

// New creates a new sprite based clip
func New(sprite *sprites.Sprite, name string, x, y, scale int) *Clip {
	frames := []*canvas.Image{}

	srcWidth, srcHeight := sprite.Width, sprite.Height
	for i := 0; i < sprite.Count; i++ {
		grid := sprite.Grid
		if grid == 0 {
			grid = sprite.Count
		}
		srcX := sprite.X + (i%grid)*(srcWidth+sprite.Gap)
		srcY := sprite.Y + (i/grid)*(srcHeight+sprite.Gap)
		srcRect := image.Rect(srcX, srcY, srcX+srcWidth, srcY+srcHeight)
		dstRect := image.Rect(0, 0, srcWidth, srcHeight)
		dst := image.NewRGBA(dstRect)
		draw.NearestNeighbor.Scale(dst, dstRect, *sprite.Image, srcRect, draw.Over, nil)
		frame := canvas.NewImageFromImage(dst)
		frame.ScaleMode = canvas.ImageScalePixels
		frames = append(frames, frame)
	}
	overlay := image.NewRGBA(frames[0].Image.Bounds())
	//blue := color.RGBA{0, 0, 255, 200}
	//draw.Draw(overlay, overlay.Bounds(), &image.Uniform{blue}, image.Point{0, 0}, draw.Src)
	clip := &Clip{
		container: container.NewStack(),
		name:      name,
		x:         x,
		y:         y,
		width:     srcWidth,
		height:    srcHeight,
		scale:     scale,
		overlay:   interactive.NewImage(canvas.NewImageFromImage(overlay)),
		frame:     0,
		frames:    frames,
	}
	for i := 0; i < len(clip.frames); i++ {
		if i == clip.frame {
			clip.frames[i].Show()
		} else {
			clip.frames[i].Hide()
		}
		clip.container.Add(clip.frames[i])
	}
	clip.container.Add(clip.overlay)
	return clip
}

// NewScaled creates a new 9 slice scaled sprite based clip
func NewScaled(sprite *sprites.Sprite, name string, x, y, width, height, scale int) *Clip {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	srcY := sprite.Y
	dstY := 0
	for h := 0; h < 3; h++ {
		srcHeight := sprite.Heights[h]
		dstHeight := sprite.Heights[h]
		if h == 1 {
			dstHeight = height - sprite.Heights[0] - sprite.Heights[2]
		}
		srcX := sprite.X
		dstX := 0
		for w := 0; w < 3; w++ {
			srcWidth := sprite.Widths[w]
			dstWidth := sprite.Widths[w]
			if w == 1 {
				dstWidth = width - sprite.Widths[0] - sprite.Widths[2]
			}

			srcRect := image.Rect(srcX, srcY, srcX+srcWidth, srcY+srcHeight)
			dstRect := image.Rect(dstX, dstY, dstX+dstWidth, dstY+dstHeight)
			draw.NearestNeighbor.Scale(dst, dstRect, *sprite.Image, srcRect, draw.Over, nil)

			srcX += srcWidth + sprite.Gap
			dstX += dstWidth
		}
		srcY += srcHeight + sprite.Gap
		dstY += dstHeight
	}
	frame0 := canvas.NewImageFromImage(dst)
	frame0.ScaleMode = canvas.ImageScalePixels
	overlay := image.NewRGBA(frame0.Image.Bounds())
	//blue := color.RGBA{0, 0, 255, 200}
	//draw.Draw(overlay, overlay.Bounds(), &image.Uniform{blue}, image.Point{0, 0}, draw.Src)
	clip := &Clip{
		container: container.NewStack(),
		name:      name,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		scale:     scale,
		overlay:   interactive.NewImage(canvas.NewImageFromImage(overlay)),
		frame:     0,
		frames:    []*canvas.Image{frame0},
	}
	clip.container.Add(frame0)
	clip.container.Add(clip.overlay)
	return clip
}

// GotoFrame goes to a frame of the clip
func (c *Clip) GotoFrame(frame int, refresh bool) {
	if c.frame != frame && frame >= 0 && frame < len(c.frames) {
		c.frame = frame
		dirty := false
		for i := 0; i < len(c.frames); i++ {
			if i == frame {
				if !c.frames[i].Visible() {
					c.frames[i].Show()
					dirty = true
				}
			} else {
				if c.frames[i].Visible() {
					c.frames[i].Hide()
					dirty = true
				}
			}
		}
		if dirty && refresh {
			c.container.Refresh()
		}
	}
}

// OnPress sets the mouse down handler
func (c *Clip) OnPress(handler func(left, right, middle, alt, control bool)) {
	c.onPress = handler
	c.overlay.OnMouseDown(func(ev *desktop.MouseEvent) {
		c.MouseDown(ev)
	})
}

// MouseDown handles the mouse down event
func (c *Clip) MouseDown(ev *desktop.MouseEvent) {
	if c.onPress != nil {
		c.onPress(ev.Button&desktop.MouseButtonPrimary > 0, ev.Button&desktop.MouseButtonSecondary > 0, ev.Button&desktop.MouseButtonTertiary > 0, false, false)
	}
}

// OnRelease sets the mouse up handler
func (c *Clip) OnRelease(handler func(left, right, middle, alt, control bool)) {
	c.onRelease = handler
	c.overlay.OnMouseUp(func(ev *desktop.MouseEvent) {
		c.MouseUp(ev)
	})
}

// MouseUp handles the mouse up event
func (c *Clip) MouseUp(ev *desktop.MouseEvent) {
	if c.onRelease != nil {
		c.onRelease(ev.Button&desktop.MouseButtonPrimary > 0, ev.Button&desktop.MouseButtonSecondary > 0, ev.Button&desktop.MouseButtonTertiary > 0, false, false)
	}
}

// OnEnter sets the enter handler
func (c *Clip) OnEnter(handler func(left, right, middle, alt, control bool)) {
	c.onEnter = handler
	c.overlay.OnMouseIn(func(ev *desktop.MouseEvent) {
		c.MouseIn(ev)
	})
}

// MouseUp handles the mouse up event
func (c *Clip) MouseIn(ev *desktop.MouseEvent) {
	if c.onEnter != nil {
		c.onEnter(ev.Button&desktop.MouseButtonPrimary > 0, ev.Button&desktop.MouseButtonSecondary > 0, ev.Button&desktop.MouseButtonTertiary > 0, false, false)
	}
}

// OnLeave sets the leave handler
func (c *Clip) OnLeave(handler func()) {
	c.onLeave = handler
	c.overlay.OnMouseOut(func() {
		c.MouseOut()
	})
}

// MouseOut handles the mouse out event
func (c *Clip) MouseOut() {
	if c.onLeave != nil {
		c.onLeave()
	}
}

// OnOver sets the mouse moved handler
func (c *Clip) OnOver(handler func(left, right, middle, alt, control bool)) {
	c.onOver = handler
	c.overlay.OnMouseMoved(func(ev *desktop.MouseEvent) {
		c.MouseMoved(ev)
	})
}

// MouseMoved handles the mouse moved event
func (c *Clip) MouseMoved(ev *desktop.MouseEvent) {
	if c.onOver != nil {
		c.onOver(ev.Button&desktop.MouseButtonPrimary > 0, ev.Button&desktop.MouseButtonSecondary > 0, ev.Button&desktop.MouseButtonTertiary > 0, false, false)
	}
}
