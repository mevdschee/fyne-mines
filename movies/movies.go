package movies

import (
	"encoding/json"
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/mevdschee/fyne-mines/clips"
	"github.com/mevdschee/fyne-mines/scenes"
	"github.com/mevdschee/fyne-mines/sprites"
)

// Movie is a set of scenes
type Movie struct {
	container    *fyne.Container
	currentScene *scenes.Scene
	scenes       map[string]*scenes.Scene
}

// New creates a new movie
func New() *Movie {
	return &Movie{
		container:    container.NewStack(),
		currentScene: nil,
		scenes:       map[string]*scenes.Scene{},
	}
}

// FromJSON creates a new movie from JSON
func FromJSON(spriteMap sprites.SpriteMap, data string, parameters map[string]interface{}) (*Movie, error) {
	sceneJSONs := []scenes.SceneJSON{}
	err := json.Unmarshal([]byte(data), &sceneJSONs)
	if err != nil {
		return nil, err
	}
	movie := Movie{
		container:    container.NewStack(),
		currentScene: &scenes.Scene{},
		scenes:       map[string]*scenes.Scene{},
	}
	for _, sceneJSON := range sceneJSONs {
		scene, err := scenes.FromJSON(spriteMap, sceneJSON, parameters)
		if err != nil {
			return nil, err
		}
		movie.Add(scene)
	}
	return &movie, nil
}

// GetContainer gets the container from the movie
func (m *Movie) GetContainer() *fyne.Container {
	return m.container
}

// SetSize set the size of the movie
func (m *Movie) SetSize(width, height int) {
	c := m.GetContainer()
	i := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 0, 0)))
	i.SetMinSize(fyne.NewSize(float32(width), float32(height)))
	c.Add(i)
}

// Add adds a scene to the movie
func (m *Movie) Add(scene *scenes.Scene) {
	m.scenes[scene.GetName()] = scene
	if len(m.scenes) == 1 {
		m.currentScene = scene
	}
	m.container.Add(scene.GetContainer())
}

// GetClip gets a clip from the movie
func (m *Movie) GetClip(scene, layer, clip string) (*clips.Clip, error) {
	return m.getClip(scene, layer, clip, 0)
}

// GetClip gets a clip from the movie
func (m *Movie) getClip(scene, layer, clip string, i int) (*clips.Clip, error) {
	if s, ok := m.scenes[scene]; ok {
		return s.GetClip(layer, clip, i)
	}
	return nil, fmt.Errorf("getClip: scene '%s' not found", scene)
}

// GetClips gets a series of clips from the movie
func (m *Movie) GetClips(scene, layer, clip string) ([]*clips.Clip, error) {
	clips := []*clips.Clip{}
	for i := 0; true; i++ {
		c, err := m.getClip(scene, layer, clip, i)
		if err != nil {
			if i == 0 {
				return clips, err
			}
			break
		}
		clips = append(clips, c)
	}
	return clips, nil
}
