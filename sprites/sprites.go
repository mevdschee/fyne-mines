package sprites

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
)

// SpriteMap is a map of sprites
type SpriteMap = map[string]*Sprite

// Sprite is the base struct for any sprite
type Sprite struct {
	Image   *image.Image `json:"-"`
	Name    string       `json:"name"`
	X       int          `json:"x"`
	Y       int          `json:"y"`
	Width   int          `json:"width,omitempty"`
	Height  int          `json:"height,omitempty"`
	Widths  [3]int       `json:"widths,omitempty"`
	Heights [3]int       `json:"heights,omitempty"`
	Count   int          `json:"count"`
	Grid    int          `json:"grid"`
	Gap     int          `json:"gap,omitempty"`
}

// NewSpriteMap creates a new sprite map
func NewSpriteMap(imageData []byte, jsondata string) (SpriteMap, error) {
	image, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}
	sprites := []*Sprite{}
	spriteMap := SpriteMap{}
	err = json.Unmarshal([]byte(jsondata), &sprites)
	if err != nil {
		return nil, err
	}
	for _, sprite := range sprites {
		sprite.Image = &image
		spriteMap[sprite.Name] = sprite
	}
	return spriteMap, nil
}
