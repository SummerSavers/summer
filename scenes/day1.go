package scenes

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Day1 struct{}

func (*Day1) Type() string { return "Day1 Scene" }

func (*Day1) Preload() {
	files := []string{
		"palm.png",
		"cloud.png",
		"title-screen-text.png",
		"Noto_Sans/NotoSans-Regular.ttf",
		"cursor.png",
		"titlebgm.wav",
	}
	if err := engo.Files.Load(files...); err != nil {
		log.Fatalf("Day1 Scene Preload. Error: %v", err)
	}
	engo.Input.RegisterButton("up", engo.KeyW)
	engo.Input.RegisterButton("down", engo.KeyS)
	engo.Input.RegisterButton("left", engo.KeyA)
	engo.Input.RegisterButton("right", engo.KeyD)
	engo.Input.RegisterButton("A", engo.KeyJ)
	engo.Input.RegisterButton("B", engo.KeyK)
}

func (*Day1) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	common.SetBackground(color.RGBA{R: 0xd8, G: 0xee, B: 0xff, A: 0xff})
}
