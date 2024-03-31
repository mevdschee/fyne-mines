package clips

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/mevdschee/fyne-mines/interactive"
	"github.com/mevdschee/fyne-mines/sprites"
	"golang.org/x/image/draw"
)

// Clip is a set of frames
type Clip struct {
	container        *fyne.Container
	name             string
	x, y             int
	width, height    int
	scale            int
	frame            int
	frames           []*interactive.Image
	onPress          func()
	onLongPress      func()
	onRelease        func()
	onReleaseOutside func()
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

// cropImage takes an image and crops it to the specified rectangle
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}
	return simg.SubImage(crop), nil
}

// New creates a new sprite based clip
func New(sprite *sprites.Sprite, name string, x, y, scale int) *Clip {
	frames := []*interactive.Image{}

	srcWidth, srcHeight := sprite.Width, sprite.Height
	for i := 0; i < sprite.Count; i++ {
		grid := sprite.Grid
		if grid == 0 {
			grid = sprite.Count
		}
		srcX := sprite.X + (i%grid)*(srcWidth+sprite.Gap)
		srcY := sprite.Y + (i/grid)*(srcHeight+sprite.Gap)
		srcRect := image.Rect(srcX, srcY, srcX+srcWidth, srcY+srcHeight)
		src, _ := cropImage(*sprite.Image, srcRect)
		dstRect := image.Rect(0, 0, srcWidth*scale, srcHeight*scale)
		dst := image.NewRGBA(dstRect)
		draw.NearestNeighbor.Scale(dst, dstRect, src, src.Bounds(), draw.Over, nil)
		frame := interactive.NewImage(canvas.NewImageFromImage(dst), fmt.Sprintf("%s: (%v,%v) x%v", name, x, y, scale))
		frames = append(frames, frame)
	}

	clip := &Clip{
		container: container.NewMax(),
		name:      name,
		x:         x,
		y:         y,
		width:     srcWidth,
		height:    srcHeight,
		scale:     scale,
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
	return clip
}

// NewScaled creates a new 9 slice scaled sprite based clip
func NewScaled(sprite *sprites.Sprite, name string, x, y, width, height, scale int) *Clip {
	dst := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))

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
			src, _ := cropImage(*sprite.Image, srcRect)
			dstRect := image.Rect(dstX*scale, dstY*scale, (dstX+dstWidth)*scale, (dstY+dstHeight)*scale)
			draw.NearestNeighbor.Scale(dst, dstRect, src, src.Bounds(), draw.Over, nil)

			srcX += srcWidth + sprite.Gap
			dstX += dstWidth
		}
		srcY += srcHeight + sprite.Gap
		dstY += dstHeight
	}
	frame0 := interactive.NewImage(canvas.NewImageFromImage(dst), fmt.Sprintf("%s: (%v,%v) x%v", name, x, y, scale))
	clip := &Clip{
		container: container.NewMax(),
		name:      name,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		scale:     scale,
		frame:     0,
		frames:    []*interactive.Image{frame0},
	}
	clip.container.Add(frame0)
	return clip
}

// Draw draws the clip
func (c *Clip) Draw(screen *canvas.Image) {
	//img := c.frames[c.frame]
	//srcWidth, srcHeight := img.Size()
	//op := &canvas.DrawImageOptions{}
	//op.GeoM.Scale(float64(c.width)/float64(srcWidth), float64(c.height)/float64(srcHeight))
	//op.GeoM.Translate(float64(c.x), float64(c.y))
	//screen.DrawImage(img, op)
}

// GotoFrame goes to a frame of the clip
func (c *Clip) GotoFrame(frame int) {
	if c.frame != frame && frame >= 0 && frame < len(c.frames) {
		c.frame = frame
		for i := 0; i < len(c.frames); i++ {
			if i == frame {
				c.frames[i].Show()
			} else {
				c.frames[i].Hide()
			}
		}
	}
}

// OnPress sets the click handler function
func (c *Clip) OnPress(handler func()) {
	c.onPress = handler
	for i := 0; i < len(c.frames); i++ {
		c.frames[i].OnMouseDown = handler
	}
}

// OnLongPress sets the click handler function
func (c *Clip) OnLongPress(handler func()) {
	c.onLongPress = handler
}

// OnRelease sets the click handler function
func (c *Clip) OnRelease(handler func()) {
	c.onRelease = handler
	for i := 0; i < len(c.frames); i++ {
		c.frames[i].OnMouseUp = handler
	}
}

// OnReleaseOutside sets the click handler function
func (c *Clip) OnReleaseOutside(handler func()) {
	c.onReleaseOutside = handler
}

// IsHovered returns whether or not the cursor is hovering the clip
func (c *Clip) IsHovered() bool {
	//cursorX, cursorY := canvas.CursorPosition()
	//cursor := image.Point{cursorX, cursorY}
	//rect := image.Rect(c.x, c.y, c.x+c.width, c.y+c.height)
	//return cursor.In(rect)
	return true
}

// Update updates the clip
func (c *Clip) Update() (err error) {
	// hover := c.IsHovered()
	// if c.onPress != nil {
	// 	if hover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	// 		c.onPress(-1)
	// 	}
	// }
	// if c.onLongPress != nil {
	// 	if hover && inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) == ebiten.MaxTPS()/2 {
	// 		c.onLongPress(-1)
	// 	}
	// }
	// if c.onRelease != nil {
	// 	if hover && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
	// 		c.onRelease(-1)
	// 	}
	// }
	// if c.onReleaseOutside != nil {
	// 	if !hover && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
	// 		c.onReleaseOutside(-1)
	// 	}
	// }
	return nil
}
