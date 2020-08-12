package scenes

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/summersavers/summer/shaders"
	"github.com/summersavers/summer/systems"
)

type TitleScene struct{}

func (*TitleScene) Type() string { return "Title Scene" }

func (*TitleScene) Preload() {
	files := []string{
		"title/palm.png",
		"title/cloud.png",
		"title/title-screen-text.png",
		"title/Noto_Sans/NotoSans-Regular.ttf",
		"title/cursor.png",
		"sounds/titlebgm.wav",
	}
	if err := engo.Files.Load(files...); err != nil {
		log.Fatalf("Title Scene Preload. Error: %v", err)
	}
	common.AddShader(shaders.CoolShader)
	engo.Input.RegisterButton("up", engo.KeyW)
	engo.Input.RegisterButton("down", engo.KeyS)
	engo.Input.RegisterButton("left", engo.KeyA)
	engo.Input.RegisterButton("right", engo.KeyD)
	engo.Input.RegisterButton("A", engo.KeyJ)
	engo.Input.RegisterButton("B", engo.KeyK)
	engo.Input.RegisterButton("X", engo.KeyL)
	engo.Input.RegisterButton("Y", engo.KeySemicolon)
	engo.Input.RegisterButton("FullScreen", engo.KeyFour)
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (*TitleScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	w.AddSystem(&systems.FullScreenSystem{})
	w.AddSystem(&systems.ExitSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	var curSys systems.CursorSystem
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	common.SetBackground(color.RGBA{R: 0xd8, G: 0xee, B: 0xff, A: 0xff})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer("sounds/titlebgm.wav")
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	cloud := sprite{BasicEntity: ecs.NewBasic()}
	cloudTex, err := common.LoadedSprite("title/cloud.png")
	if err != nil {
		log.Fatalf("Title Scene Setup. cloud.png texture was not found. Error was: %v", err)
	}
	cloud.RenderComponent.Drawable = cloudTex
	cloud.RenderComponent.Scale.Set(10, 10)
	cloud.SpaceComponent.Position.Set(40, 0)
	w.AddEntity(&cloud)

	tree := sprite{BasicEntity: ecs.NewBasic()}
	treeTex, err := common.LoadedSprite("title/palm.png")
	if err != nil {
		log.Fatalf("Title Scene Setup. palm.png texture was not found. Error was: %v", err)
	}
	tree.RenderComponent.Drawable = treeTex
	tree.RenderComponent.Scale = engo.Point{X: 6, Y: 6}
	tree.RenderComponent.SetZIndex(1)
	tree.SpaceComponent.Position.Set(22, 165)
	w.AddEntity(&tree)

	title := sprite{BasicEntity: ecs.NewBasic()}
	titleTex, err := common.LoadedSprite("title/title-screen-text.png")
	title.Drawable = titleTex
	title.RenderComponent.SetZIndex(1)
	title.SpaceComponent.Position.Set(170, 70)
	w.AddEntity(&title)

	selFont := &common.Font{
		Size: 20,
		FG:   color.Black,
		URL:  "title/Noto_Sans/NotoSans-Regular.ttf",
	}
	selFont.CreatePreloaded()

	startText := selection{BasicEntity: ecs.NewBasic()}
	startText.Drawable = common.Text{
		Text: "start",
		Font: selFont,
	}
	startText.SetZIndex(1)
	startText.SetCenter(engo.Point{X: 306, Y: 150})
	startText.CursorComponent.ACallback = func() {
		engo.SetSceneByName("intro battle", true)
	}
	w.AddEntity(&startText)

	fsText := sprite{BasicEntity: ecs.NewBasic()}
	fsText.Drawable = common.Text{
		Text: "press 4 to enable full screen!",
		Font: selFont,
	}
	fsText.SetZIndex(1)
	fsText.SetCenter(engo.Point{X: 250, Y: 180})
	fsText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&fsText)
}
