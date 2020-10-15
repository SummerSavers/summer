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

type IntroScene struct{}

func (*IntroScene) Type() string { return "intro battle" }

var files = []string{
	"combat/BloodmouthGhost/font.ttf",
	"combat/BloodmouthGhost/dots.png",
	"combat/BloodmouthGhost/Logs.png",
	"combat/BloodmouthGhost/BloodmouthGhost.png",
	"combat/Griff/griffcard.png",
	"combat/Kai/kaicard.png",
	"combat/Levi/levicard.png",
	"combat/Wy/wycard.png",
	"combat/Griff/font.ttf",
	"combat/Kai/font.ttf",
	"combat/Levi/font.ttf",
	"combat/Wy/font.ttf",
	"combat/Griff/GriffBash.png",
	"combat/Griff/GriffBlock.png",
	"combat/Griff/GriffEmpty.png",
	"combat/Griff/GriffStun.png",
	"combat/Kai/KaiAxe.png",
	"combat/Kai/KaiEmpty.png",
	"combat/Kai/KaiEnrage.png",
	"combat/Kai/KaiFlare.png",
	"combat/Levi/LeviEmpty.png",
	"combat/Levi/LeviFreeze.png",
	"combat/Levi/LeviLightning.png",
	"combat/Levi/LeviStaff.png",
	"combat/Cent/CentEmpty.png",
	"combat/Cent/CentFlare.png",
	"combat/Cent/CentHeal.png",
	"combat/Cent/CentSlingshot.png",
	"combat/Wy/WyEmpty.png",
	"combat/Wy/WyEvade.png",
	"combat/Wy/WySteal.png",
	"combat/Wy/WySword.png",
	"combat/Cent/cencard.png",
	"combat/Cent/font.ttf",
	"combat/BattleBoxes.png",
	"sounds/bloodmouthghost/log.wav",
	"sounds/griff/log.wav",
	"sounds/cent/log.wav",
	"sounds/kai/log.wav",
	"sounds/levi/log.wav",
	"sounds/wy/log.wav",
	"sounds/bloodmouthghost/bg.wav",
}

func (*IntroScene) Preload() {
	if err := engo.Files.Load(files...); err != nil {
		log.Fatalf("Intro Scene Preload. Error: %v", err)
	}
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

func (*IntroScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	// < SYSTEMS
	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	w.AddSystem(&systems.CombatLogSystem{
		BackgroundURL: files[2],
		DotURL:        files[1],
		FontURL:       files[0],
	})

	var baddiephaseable *systems.BaddiePhaseAble
	var notbaddiephaseable *systems.NotBaddiePhaseAble
	w.AddSystemInterface(&systems.BaddiePhaseSystem{}, baddiephaseable, notbaddiephaseable)

	var playerselectable *systems.PlayerSelectAble
	var notplayerselectable *systems.NotPlayerSelectAble
	w.AddSystemInterface(&systems.PlayerSelectSystem{PlayerCount: 5}, playerselectable, notplayerselectable)

	w.AddSystem(&systems.ExitSystem{})
	w.AddSystem(&systems.TargetSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	w.AddSystemInterface(&systems.CursorSystem{}, cursorable, notcursorable)

	var battleboxable *systems.BattleboxAble
	var notbattleboxable *systems.NotBattleboxAble
	w.AddSystemInterface(&systems.BattleboxSystem{}, battleboxable, notbattleboxable)

	var characterbarable *systems.CharacterBarAble
	var baddiebarable *systems.BaddieBarAble
	var notbarable *systems.NotBarAble
	w.AddSystemInterface(&systems.BarSystem{}, []interface{}{characterbarable, baddiebarable}, notbarable)
	// SYSTEMS />

	// < BACKGROUNDS
	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = shaders.ShaderDrawable{}
	bg.SetShader(shaders.CoolShader)
	w.AddEntity(&bg)
	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(files[41])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)
	// BACKGROUNDS />

	// < BloodmouthGhost
	//    < sounds
	bmgLogClip := audio{BasicEntity: ecs.NewBasic()}
	bmgLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[35])
	bmgLogClip.Player.SetVolume(0.2)
	w.AddEntity(&bmgLogClip)
	//    sounds />
	//    < sprite
	bmg := baddie{BasicEntity: ecs.NewBasic()}
	bmg.SetZIndex(2)
	bmg.SetCenter(engo.Point{X: 320, Y: 130})
	bmg.BaddieComponent = systems.BaddieComponent{
		Name:         "Blood-Mouthed Ghost!",
		Phases:       make(map[string]systems.Phase),
		CurrentPhase: "Start",
		Clip:         bmgLogClip.Player,
		HP:           200,
		MaxHP:        200,
		Spritesheet:  common.NewSpritesheetWithBorderFromFile(files[3], 64, 64, 1, 1),
	}
	bmg.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[0],
	}
	bmg.Font.CreatePreloaded()
	bmg.Drawable = bmg.Spritesheet.Drawable(0)
	bmg.Width = bmg.Drawable.Width()
	bmg.Height = bmg.Drawable.Height()
	w.AddEntity(&bmg)

	for _, sys := range w.Systems() {
		switch system := sys.(type) {
		case *systems.TargetSystem:
			system.Add(&bmg.BasicEntity, &bmg.RenderComponent, &bmg.BaddieComponent)
		}
	}
	//    sprite />
	// BloodmouthGhost />

	//    < battlebox sprite sheet
	battleboxSpritesheet := common.NewSpritesheetWithBorderFromFile(files[34], 600, 144, 1, 1)
	//    battlebox sprite sheet />

	// < Kids
	// <Griff
	griff := character{BasicEntity: ecs.NewBasic()}
	griff.Name = "Griff"
	griff.Abilities = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Ability A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability D",
			Desc: "!!!",
		},
	}
	griff.Items = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Item A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item D",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item E",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item F",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item G",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item H",
			Desc: "...",
		},
	}
	griff.Acts = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Act A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act C",
			Desc: "...",
		},
	}
	griff.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[8],
	}
	griff.Font.CreatePreloaded()
	griffLogClip := audio{BasicEntity: ecs.NewBasic()}
	griffLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[36])
	griffLogClip.Player.SetVolume(0.2)
	w.AddEntity(&griffLogClip)
	griff.Clip = griffLogClip.Player
	griff.Card, _ = common.LoadedSprite(files[4])
	griff.BattleBox = battleboxSpritesheet.Drawable(4)
	griff.AIcon, _ = common.LoadedSprite(files[12])
	griff.BIcon, _ = common.LoadedSprite(files[15])
	griff.XIcon, _ = common.LoadedSprite(files[13])
	griff.YIcon, _ = common.LoadedSprite(files[14])
	griff.HP = 100
	griff.MaxHP = 100
	griff.MP = 100
	griff.MaxMP = 100
	griff.CardSelected = true
	w.AddEntity(&griff)
	// Griff />

	// < Kai
	kai := character{BasicEntity: ecs.NewBasic()}
	kai.Name = "Kai"
	kai.Abilities = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Ability A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability D",
			Desc: "!!!",
		},
	}
	kai.Items = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Item A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item D",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item E",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item F",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item G",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item H",
			Desc: "...",
		},
	}
	kai.Acts = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Act A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act C",
			Desc: "...",
		},
	}
	kai.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[9],
	}
	kai.Font.CreatePreloaded()
	kaiLogClip := audio{BasicEntity: ecs.NewBasic()}
	kaiLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[38])
	kaiLogClip.Player.SetVolume(0.2)
	w.AddEntity(&kaiLogClip)
	kai.Clip = kaiLogClip.Player
	kai.Card, _ = common.LoadedSprite(files[5])
	kai.BattleBox = battleboxSpritesheet.Drawable(0)
	kai.AIcon, _ = common.LoadedSprite(files[16])
	kai.BIcon, _ = common.LoadedSprite(files[18])
	kai.XIcon, _ = common.LoadedSprite(files[19])
	kai.YIcon, _ = common.LoadedSprite(files[17])
	kai.HP = 100
	kai.MaxHP = 100
	kai.MP = 100
	kai.MaxMP = 100
	w.AddEntity(&kai)
	// Kai />
	// < Cent
	cent := character{BasicEntity: ecs.NewBasic()}
	cent.Name = "Cent"
	cent.Abilities = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Ability A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability D",
			Desc: "!!!",
		},
	}
	cent.Items = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Item A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item D",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item E",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item F",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item G",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item H",
			Desc: "...",
		},
	}
	cent.Acts = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Act A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act C",
			Desc: "...",
		},
	}
	cent.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[33],
	}
	cent.Font.CreatePreloaded()
	centLogClip := audio{BasicEntity: ecs.NewBasic()}
	centLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[38])
	centLogClip.Player.SetVolume(0.2)
	w.AddEntity(&centLogClip)
	cent.Clip = centLogClip.Player
	cent.Card, _ = common.LoadedSprite(files[32])
	cent.BattleBox = battleboxSpritesheet.Drawable(1)
	cent.AIcon, _ = common.LoadedSprite(files[27])
	cent.BIcon, _ = common.LoadedSprite(files[25])
	cent.XIcon, _ = common.LoadedSprite(files[26])
	cent.YIcon, _ = common.LoadedSprite(files[24])
	cent.HP = 100
	cent.MaxHP = 100
	cent.MP = 100
	cent.MaxMP = 100
	w.AddEntity(&cent)
	// Cent />
	// <Levi
	levi := character{BasicEntity: ecs.NewBasic()}
	levi.Name = "Levi"
	levi.Abilities = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Ability A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability D",
			Desc: "!!!",
		},
	}
	levi.Items = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Item A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item D",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item E",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item F",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item G",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item H",
			Desc: "...",
		},
	}
	levi.Acts = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Act A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act C",
			Desc: "...",
		},
	}
	levi.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[10],
	}
	levi.Font.CreatePreloaded()
	leviLogClip := audio{BasicEntity: ecs.NewBasic()}
	leviLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[39])
	leviLogClip.Player.SetVolume(0.2)
	w.AddEntity(&leviLogClip)
	levi.Clip = leviLogClip.Player
	levi.Card, _ = common.LoadedSprite(files[6])
	levi.BattleBox = battleboxSpritesheet.Drawable(2)
	levi.AIcon, _ = common.LoadedSprite(files[23])
	levi.BIcon, _ = common.LoadedSprite(files[21])
	levi.XIcon, _ = common.LoadedSprite(files[22])
	levi.YIcon, _ = common.LoadedSprite(files[20])
	levi.HP = 100
	levi.MaxHP = 100
	levi.MP = 100
	levi.MaxMP = 100
	w.AddEntity(&levi)
	// Levi />
	// <Wy
	wy := character{BasicEntity: ecs.NewBasic()}
	wy.Name = "Wy"
	wy.Abilities = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Ability A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Ability D",
			Desc: "!!!",
		},
	}
	wy.Items = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Item A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item C",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item D",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item E",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item F",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item G",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Item H",
			Desc: "...",
		},
	}
	wy.Acts = []systems.BattleBoxText{
		systems.BattleBoxText{
			Name: "Act A",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act B",
			Desc: "...",
		},
		systems.BattleBoxText{
			Name: "Act C",
			Desc: "...",
		},
	}
	wy.Font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[11],
	}
	wy.Font.CreatePreloaded()
	wyLogClip := audio{BasicEntity: ecs.NewBasic()}
	wyLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[40])
	wyLogClip.Player.SetVolume(0.2)
	w.AddEntity(&wyLogClip)
	wy.Clip = wyLogClip.Player
	wy.Card, _ = common.LoadedSprite(files[7])
	wy.BattleBox = battleboxSpritesheet.Drawable(3)
	wy.AIcon, _ = common.LoadedSprite(files[31])
	wy.BIcon, _ = common.LoadedSprite(files[29])
	wy.XIcon, _ = common.LoadedSprite(files[30])
	wy.YIcon, _ = common.LoadedSprite(files[28])
	wy.HP = 100
	wy.MaxHP = 100
	wy.MP = 100
	wy.MaxMP = 100
	w.AddEntity(&wy)
	// Wy />

	// < Send a BMG Message
	engo.Mailbox.Dispatch(systems.CombatLogMessage{
		Msg:  "A blood-mouthed ghost has appearerated!",
		Fnt:  bmg.Font,
		Clip: bmgLogClip.Player,
	})
	// Send a BMG Message />
}
