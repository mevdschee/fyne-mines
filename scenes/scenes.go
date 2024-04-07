package scenes

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/mevdschee/fyne-mines/clips"
	"github.com/mevdschee/fyne-mines/layers"
	"github.com/mevdschee/fyne-mines/sprites"
)

// Scene is a set of layers
type Scene struct {
	container *fyne.Container
	name      string
	layers    map[string]*layers.Layer
	order     []string
}

// SceneJSON is a set of layers in JSON
type SceneJSON struct {
	Name   string
	Layers []layers.LayerJSON
}

// GetName gets the name of the scene
func (s *Scene) GetName() string {
	return s.name
}

// GetLayers gets the layers of the scene
func (s *Scene) GetLayers() map[string]*layers.Layer {
	return s.layers
}

// New creates a new scene
func New(name string) *Scene {
	return &Scene{
		container: container.NewStack(),
		name:      name,
		layers:    map[string]*layers.Layer{},
		order:     []string{},
	}
}

// FromJSON creates a new scene from JSON
func FromJSON(spriteMap sprites.SpriteMap, sceneJSON SceneJSON, parameters map[string]interface{}) (*Scene, error) {
	scene := Scene{
		container: container.NewStack(),
		name:      sceneJSON.Name,
		layers:    map[string]*layers.Layer{},
		order:     []string{},
	}
	for _, layerJSON := range sceneJSON.Layers {
		layer, err := layers.FromJSON(spriteMap, layerJSON, parameters)
		if err != nil {
			return nil, err
		}
		scene.Add(layer)
	}
	return &scene, nil
}

// GetContainer gets the container from the scene
func (s *Scene) GetContainer() *fyne.Container {
	return s.container
}

// Add adds a layers to the scene
func (s *Scene) Add(layer *layers.Layer) {
	name := layer.GetName()
	s.layers[name] = layer
	s.order = append(s.order, name)
	s.container.Add(layer.GetContainer())
}

// GetClip gets a clip from the scene
func (s *Scene) GetClip(layer, clip string, i int) (*clips.Clip, error) {
	if l, ok := s.layers[layer]; ok {
		return l.GetClip(clip, i)
	}
	return nil, fmt.Errorf("GetClip: layer '%s' not found", layer)
}
