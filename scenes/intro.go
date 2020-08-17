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
	w.AddSystemInterface(&systems.PlayerSelectSystem{PlayerCount: 1}, playerselectable, notplayerselectable)

	w.AddSystem(&systems.ExitSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	w.AddSystemInterface(&systems.CursorSystem{}, cursorable, notcursorable)

	var battleboxable *systems.BattleboxAble
	var notbattleboxable *systems.NotBattleboxAble
	w.AddSystemInterface(&systems.BattleboxSystem{}, battleboxable, notbattleboxable)
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
	bmgFont := &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  files[0],
	}
	bmgFont.CreatePreloaded()
	//    < sprite sheet
	bmgSpritesheet := common.NewSpritesheetWithBorderFromFile(files[3], 64, 64, 1, 1)
	//    sprite sheet />
	//    < sounds
	bmgLogClip := audio{BasicEntity: ecs.NewBasic()}
	bmgLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[35])
	bmgLogClip.Player.SetVolume(0.2)
	w.AddEntity(&bmgLogClip)
	//    sounds />
	//    < sprite
	bmgSprite := baddieSprite{BasicEntity: ecs.NewBasic()}
	bmgSprite.Drawable = bmgSpritesheet.Drawable(0)
	bmgSprite.Width = bmgSprite.Drawable.Width()
	bmgSprite.Height = bmgSprite.Drawable.Height()
	bmgSprite.SetZIndex(2)
	bmgSprite.SetCenter(engo.Point{X: 320, Y: 130})
	w.AddEntity(&bmgSprite)
	//    sprite />
	// BloodmouthGhost />

	//    < battlebox sprite sheet
	battleboxSpritesheet := common.NewSpritesheetWithBorderFromFile(files[34], 600, 144, 1, 1)
	//    battlebox sprite sheet />

	// // < Kids
	// < Griff
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

	// < Send a BMG Message
	engo.Mailbox.Dispatch(systems.CombatLogMessage{
		Msg:  "A blood-mouthed ghost has appearerated!",
		Fnt:  bmgFont,
		Clip: bmgLogClip.Player,
	})
	// Send a BMG Message />

	// // Kids />
	// // < GRIFF
	// //    <Load font
	// griffFont := &common.Font{
	// 	Size: 64,
	// 	FG:   color.Black,
	// 	URL:  files[8],
	// }
	// griffFont.CreatePreloaded()
	// //    Load font />
	// //    < Sounds
	// griffLogClip := audio{BasicEntity: ecs.NewBasic()}
	// griffLogClip.AudioComponent.Player, _ = common.LoadedPlayer(files[36])
	// griffLogClip.Player.SetVolume(0.2)
	// w.AddEntity(&griffLogClip)
	// //    Sounds />
	// //    < CARD
	// griffCard := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.Drawable, _ = common.LoadedSprite(files[4])
	// griffCard.SetZIndex(3)
	// griffCard.Position = engo.Point{X: 21, Y: 192}
	// griffCard.Selected = true
	// griffCard.Card = true
	// w.AddEntity(&griffCard)
	// //    CARD />
	// //        <NAME
	// griffName := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffName.BasicEntity)
	// griffName.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Griff",
	// }
	// griffName.SetZIndex(4)
	// griffName.Position = engo.Point{X: 27, Y: 197}
	// griffName.Scale = engo.Point{X: 0.2, Y: 0.2}
	// w.AddEntity(&griffName)
	// //        NAME />
	// //        <HPBAR
	// griffHPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffHPBar.BasicEntity)
	// griffHPBar.Drawable = common.Rectangle{}
	// griffHPBar.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// griffHPBar.SetZIndex(4)
	// griffHPBar.Width = 83
	// griffHPBar.Height = 13
	// griffHPBar.Position = engo.Point{X: 29, Y: 224}
	// w.AddEntity(&griffHPBar)
	// //        HPBAR />
	// //        <MPBAR
	// griffMPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffMPBar.BasicEntity)
	// griffMPBar.Drawable = common.Rectangle{}
	// griffMPBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	// griffMPBar.SetZIndex(4)
	// griffMPBar.Width = 83
	// griffMPBar.Height = 13
	// griffMPBar.Position = engo.Point{X: 29, Y: 246}
	// w.AddEntity(&griffMPBar)
	// //        MPBAR />
	// //        <CastBar
	// griffCBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffCBar.BasicEntity)
	// griffCBar.Drawable = common.Rectangle{}
	// griffCBar.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	// griffCBar.SetZIndex(4)
	// griffCBar.Width = 83
	// griffCBar.Height = 13
	// griffCBar.Position = engo.Point{X: 29, Y: 268}
	// w.AddEntity(&griffCBar)
	// //        CastBar />
	// //        <MenuBG
	// griffBattleBoxBG := sprite{BasicEntity: ecs.NewBasic()}
	// griffBattleBoxBG.Drawable = battleboxSpritesheet.Drawable(4)
	// griffBattleBoxBG.SetZIndex(5)
	// griffBattleBoxBG.Position.Set(20, 150)
	// griffBattleBoxBG.Color = color.RGBA{0xff, 0xff, 0xff, 0xe5}
	// griffBattleBoxBG.Hidden = true
	// w.AddEntity(&griffBattleBoxBG)
	// //        MenuBG />
	// //        <Menu Special
	// //            < Bash
	// griffBash := selection{BasicEntity: ecs.NewBasic()}
	// griffBash.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Bash",
	// }
	// griffBash.SetZIndex(6)
	// griffBash.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffBash.Position = engo.Point{X: 55, Y: 175}
	// griffBash.Hidden = true
	// w.AddEntity(&griffBash)
	// //            Bash />
	// //            < Shield
	// griffShield := selection{BasicEntity: ecs.NewBasic()}
	// griffShield.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Shield",
	// }
	// griffShield.SetZIndex(6)
	// griffShield.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffShield.Position = engo.Point{X: 355, Y: 175}
	// griffShield.Hidden = true
	// w.AddEntity(&griffShield)
	// //            Shield />
	// //            < Cover
	// griffCover := selection{BasicEntity: ecs.NewBasic()}
	// griffCover.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Cover",
	// }
	// griffCover.SetZIndex(6)
	// griffCover.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffCover.Position = engo.Point{X: 55, Y: 235}
	// griffCover.Hidden = true
	// w.AddEntity(&griffCover)
	// //            Cover />
	// //            < Taunt
	// griffTaunt := selection{BasicEntity: ecs.NewBasic()}
	// griffTaunt.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Taunt",
	// }
	// griffTaunt.SetZIndex(6)
	// griffTaunt.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffTaunt.Position = engo.Point{X: 355, Y: 235}
	// griffTaunt.Hidden = true
	// w.AddEntity(&griffTaunt)
	// //            Taunt />
	// //        Menu Special />
	// //        < Menu Items
	// //            < Item 1
	// griffItem1 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem1.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Candy Bar",
	// }
	// griffItem1.SetZIndex(6)
	// griffItem1.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem1.Position = engo.Point{X: 55, Y: 175}
	// // griffItem1.CursorComponent.ACallback = func(s *systems.CursorSystem) {
	// // 	engo.Mailbox.Dispatch(systems.CombatLogMessage{
	// // 		Msg:  "candy yum yum",
	// // 		Fnt:  griffFont,
	// // 		Clip: griffLogClip.Player,
	// // 	})
	// // }
	// griffItem1.Hidden = true
	// w.AddEntity(&griffItem1)
	// //            Item 1 />
	// //            < Item 2
	// griffItem2 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem2.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Stick",
	// }
	// griffItem2.SetZIndex(6)
	// griffItem2.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem2.Position = engo.Point{X: 355, Y: 175}
	// griffItem2.Hidden = true
	// w.AddEntity(&griffItem2)
	// //            Item 2 />
	// //            < Item 3
	// griffItem3 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem3.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Salt",
	// }
	// griffItem3.SetZIndex(6)
	// griffItem3.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem3.Position = engo.Point{X: 55, Y: 235}
	// griffItem3.Hidden = true
	// w.AddEntity(&griffItem3)
	// //            Item 3 />
	// //            < Item 4
	// griffItem4 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem4.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Saltwater taffy",
	// }
	// griffItem4.SetZIndex(6)
	// griffItem4.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem4.Position = engo.Point{X: 355, Y: 235}
	// griffItem4.Hidden = true
	// w.AddEntity(&griffItem4)
	// //            Item 4 />
	// //            < Item 5
	// griffItem5 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem5.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Sea Salt",
	// }
	// griffItem5.SetZIndex(6)
	// griffItem5.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem5.Position = engo.Point{X: 55, Y: 175}
	// griffItem5.Hidden = true
	// w.AddEntity(&griffItem5)
	// //            Item 5 />
	// //            < Item 6
	// griffItem6 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem6.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Sportsade",
	// }
	// griffItem6.SetZIndex(6)
	// griffItem6.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem6.Position = engo.Point{X: 355, Y: 175}
	// griffItem6.Hidden = true
	// w.AddEntity(&griffItem6)
	// //            Item 6 />
	// //            < Item 7
	// griffItem7 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem7.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Salt",
	// }
	// griffItem7.SetZIndex(6)
	// griffItem7.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem7.Position = engo.Point{X: 55, Y: 235}
	// griffItem7.Hidden = true
	// w.AddEntity(&griffItem7)
	// //            Item 7 />
	// //            < Item 8
	// griffItem8 := selection{BasicEntity: ecs.NewBasic()}
	// griffItem8.Drawable = common.Text{
	// 	Font: griffFont,
	// 	Text: "Sportsade",
	// }
	// griffItem8.SetZIndex(6)
	// griffItem8.Scale = engo.Point{X: 0.4, Y: 0.4}
	// griffItem8.Position = engo.Point{X: 355, Y: 235}
	// griffItem8.Hidden = true
	// w.AddEntity(&griffItem8)
	// //            Item 8 />
	// //        Menu Items />
	// //        <Menu Act
	//
	// //        Menu Act />
	// //        <ACTION A
	// griffA := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffA.BasicEntity)
	// griffA.Drawable, _ = common.LoadedSprite(files[])
	// griffA.Width = griffA.Drawable.Width()
	// griffA.Height = griffA.Drawable.Height()
	// griffA.SetZIndex(4)
	// griffA.SetCenter(engo.Point{X: 45, Y: 294})
	// w.AddEntity(&griffA)
	// //        ACTION A />
	// //        <ACTION B
	// griffB := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffB.BasicEntity)
	// griffB.Drawable, _ = common.LoadedSprite(files[])
	// griffB.Width = griffB.Drawable.Width()
	// griffB.Height = griffB.Drawable.Height()
	// griffB.SetZIndex(4)
	// griffB.SetCenter(engo.Point{X: 96, Y: 294})
	// w.AddEntity(&griffB)
	// //        ACTION B />
	// //        <ACTION X
	// griffX := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffX.BasicEntity)
	// griffX.Drawable, _ = common.LoadedSprite(files[])
	// griffX.Width = griffX.Drawable.Width()
	// griffX.Height = griffX.Drawable.Height()
	// griffX.SetZIndex(4)
	// griffX.SetCenter(engo.Point{X: 45, Y: 318})
	// w.AddEntity(&griffX)
	// //        ACTION X />
	// //        <ACTION Y
	// griffY := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// griffCard.AppendChild(&griffY.BasicEntity)
	// griffY.Drawable, _ = common.LoadedSprite(files[])
	// griffY.Width = griffX.Drawable.Width()
	// griffY.Height = griffX.Drawable.Height()
	// griffY.SetZIndex(4)
	// griffY.SetCenter(engo.Point{X: 96, Y: 318})
	// w.AddEntity(&griffY)
	// //        ACTION Y />
	// // GRIFF />
	//
	// // < KAI
	// //    <Load font
	// kaiFont := &common.Font{
	// 	Size: 64,
	// 	FG:   color.Black,
	// 	URL:  files[9],
	// }
	// kaiFont.CreatePreloaded()
	// //    Load font />
	// //    < CARD
	// kaiCard := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.Drawable, _ = common.LoadedSprite(files[5])
	// kaiCard.SetZIndex(3)
	// kaiCard.Position = engo.Point{X: 148, Y: 217}
	// kaiCard.Card = true
	// w.AddEntity(&kaiCard)
	// //    CARD />
	// //        <NAME
	// kaiName := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiName.BasicEntity)
	// kaiName.Drawable = common.Text{
	// 	Font: kaiFont,
	// 	Text: "Kai",
	// }
	// kaiName.SetZIndex(4)
	// kaiName.Position = engo.Point{X: 154, Y: 222}
	// kaiName.Scale = engo.Point{X: 0.2, Y: 0.2}
	// w.AddEntity(&kaiName)
	// //        NAME />
	// //        <HPBAR
	// kaiHPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiHPBar.BasicEntity)
	// kaiHPBar.Drawable = common.Rectangle{}
	// kaiHPBar.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// kaiHPBar.SetZIndex(4)
	// kaiHPBar.Width = 83
	// kaiHPBar.Height = 13
	// kaiHPBar.Position = engo.Point{X: 156, Y: 249}
	// w.AddEntity(&kaiHPBar)
	// //        HPBAR />
	// //        <MPBAR
	// kaiMPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiMPBar.BasicEntity)
	// kaiMPBar.Drawable = common.Rectangle{}
	// kaiMPBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	// kaiMPBar.SetZIndex(4)
	// kaiMPBar.Width = 83
	// kaiMPBar.Height = 13
	// kaiMPBar.Position = engo.Point{X: 156, Y: 271}
	// w.AddEntity(&kaiMPBar)
	// //        MPBAR />
	// //        <CastBar
	// kaiCBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiCBar.BasicEntity)
	// kaiCBar.Drawable = common.Rectangle{}
	// kaiCBar.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	// kaiCBar.SetZIndex(4)
	// kaiCBar.Width = 83
	// kaiCBar.Height = 13
	// kaiCBar.Position = engo.Point{X: 156, Y: 293}
	// w.AddEntity(&kaiCBar)
	// //        CastBar />
	// //        <ACTION A
	// kaiA := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiA.BasicEntity)
	// kaiA.Drawable, _ = common.LoadedSprite(files[16])
	// kaiA.Width = kaiA.Drawable.Width()
	// kaiA.Height = kaiA.Drawable.Height()
	// kaiA.SetZIndex(4)
	// kaiA.SetCenter(engo.Point{X: 172, Y: 319})
	// w.AddEntity(&kaiA)
	// //        ACTION A />
	// //        <ACTION B
	// kaiB := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiB.BasicEntity)
	// kaiB.Drawable, _ = common.LoadedSprite(files[18])
	// kaiB.Width = kaiB.Drawable.Width()
	// kaiB.Height = kaiB.Drawable.Height()
	// kaiB.SetZIndex(4)
	// kaiB.SetCenter(engo.Point{X: 223, Y: 319})
	// w.AddEntity(&kaiB)
	// //        ACTION B />
	// //        <ACTION X
	// kaiX := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiX.BasicEntity)
	// kaiX.Drawable, _ = common.LoadedSprite(files[19])
	// kaiX.Width = kaiX.Drawable.Width()
	// kaiX.Height = kaiX.Drawable.Height()
	// kaiX.SetZIndex(4)
	// kaiX.SetCenter(engo.Point{X: 172, Y: 343})
	// w.AddEntity(&kaiX)
	// //        ACTION X />
	// //        <ACTION Y
	// kaiY := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// kaiCard.AppendChild(&kaiY.BasicEntity)
	// kaiY.Drawable, _ = common.LoadedSprite(files[17])
	// kaiY.Width = kaiY.Drawable.Width()
	// kaiY.Height = kaiY.Drawable.Height()
	// kaiY.SetZIndex(4)
	// kaiY.SetCenter(engo.Point{X: 223, Y: 343})
	// w.AddEntity(&kaiY)
	// //        ACTION Y />
	// // KAI />
	//
	// // < CENT
	// //    <Load font
	// centFont := &common.Font{
	// 	Size: 64,
	// 	FG:   color.Black,
	// 	URL:  files[33],
	// }
	// centFont.CreatePreloaded()
	// //    Load font />
	// //    < CARD
	// centCard := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.Card = true
	// centCard.Drawable, _ = common.LoadedSprite(files[32])
	// centCard.SetZIndex(3)
	// centCard.Position = engo.Point{X: 275, Y: 217}
	// w.AddEntity(&centCard)
	// //    CARD />
	// //        <NAME
	// centName := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centName.BasicEntity)
	// centName.Drawable = common.Text{
	// 	Font: centFont,
	// 	Text: "Cent",
	// }
	// centName.SetZIndex(4)
	// centName.Position = engo.Point{X: 281, Y: 222}
	// centName.Scale = engo.Point{X: 0.25, Y: 0.25}
	// w.AddEntity(&centName)
	// //        NAME />
	// //        <HPBAR
	// centHPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centHPBar.BasicEntity)
	// centHPBar.Drawable = common.Rectangle{}
	// centHPBar.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// centHPBar.SetZIndex(4)
	// centHPBar.Width = 83
	// centHPBar.Height = 13
	// centHPBar.Position = engo.Point{X: 283, Y: 249}
	// w.AddEntity(&centHPBar)
	// //        HPBAR />
	// //        <MPBAR
	// centMPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centMPBar.BasicEntity)
	// centMPBar.Drawable = common.Rectangle{}
	// centMPBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	// centMPBar.SetZIndex(4)
	// centMPBar.Width = 83
	// centMPBar.Height = 13
	// centMPBar.Position = engo.Point{X: 283, Y: 271}
	// w.AddEntity(&centMPBar)
	// //        MPBAR />
	// //        <CastBar
	// centCBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centCBar.BasicEntity)
	// centCBar.Drawable = common.Rectangle{}
	// centCBar.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	// centCBar.SetZIndex(4)
	// centCBar.Width = 83
	// centCBar.Height = 13
	// centCBar.Position = engo.Point{X: 283, Y: 293}
	// w.AddEntity(&centCBar)
	// //        CastBar />
	// //        <ACTION A
	// centA := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centA.BasicEntity)
	// centA.Drawable, _ = common.LoadedSprite(files[27])
	// centA.Width = centA.Drawable.Width()
	// centA.Height = centA.Drawable.Height()
	// centA.SetZIndex(4)
	// centA.SetCenter(engo.Point{X: 299, Y: 319})
	// w.AddEntity(&centA)
	// //        ACTION A />
	// //        <ACTION B
	// centB := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centB.BasicEntity)
	// centB.Drawable, _ = common.LoadedSprite(files[25])
	// centB.Width = centB.Drawable.Width()
	// centB.Height = centB.Drawable.Height()
	// centB.SetZIndex(4)
	// centB.SetCenter(engo.Point{X: 350, Y: 319})
	// w.AddEntity(&centB)
	// //        ACTION B />
	// //        <ACTION X
	// centX := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centX.BasicEntity)
	// centX.Drawable, _ = common.LoadedSprite(files[26])
	// centX.Width = centX.Drawable.Width()
	// centX.Height = centX.Drawable.Height()
	// centX.SetZIndex(4)
	// centX.SetCenter(engo.Point{X: 299, Y: 343})
	// w.AddEntity(&centX)
	// //        ACTION X />
	// //        <ACTION Y
	// centY := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// centCard.AppendChild(&centY.BasicEntity)
	// centY.Drawable, _ = common.LoadedSprite(files[24])
	// centY.Width = centY.Drawable.Width()
	// centY.Height = centY.Drawable.Height()
	// centY.SetZIndex(4)
	// centY.SetCenter(engo.Point{X: 350, Y: 343})
	// w.AddEntity(&centY)
	// //        ACTION Y />
	// // CENT />
	//
	// // <LEVI
	// //    <Load font
	// leviFont := &common.Font{
	// 	Size: 64,
	// 	FG:   color.Black,
	// 	URL:  files[10],
	// }
	// leviFont.CreatePreloaded()
	// //    Load font />
	// //    < CARD
	// leviCard := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.Card = true
	// leviCard.Drawable, _ = common.LoadedSprite(files[6])
	// leviCard.SetZIndex(3)
	// leviCard.Position = engo.Point{X: 402, Y: 217}
	// w.AddEntity(&leviCard)
	// //    CARD />
	// //        <NAME
	// leviName := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviName.BasicEntity)
	// leviName.Drawable = common.Text{
	// 	Font: leviFont,
	// 	Text: "Levi",
	// }
	// leviName.SetZIndex(4)
	// leviName.Position = engo.Point{X: 408, Y: 222}
	// leviName.Scale = engo.Point{X: 0.25, Y: 0.25}
	// w.AddEntity(&leviName)
	// //        NAME />
	// //        <HPBAR
	// leviHPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviHPBar.BasicEntity)
	// leviHPBar.Drawable = common.Rectangle{}
	// leviHPBar.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// leviHPBar.SetZIndex(4)
	// leviHPBar.Width = 83
	// leviHPBar.Height = 13
	// leviHPBar.Position = engo.Point{X: 410, Y: 249}
	// w.AddEntity(&leviHPBar)
	// //        HPBAR />
	// //        <MPBAR
	// leviMPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviMPBar.BasicEntity)
	// leviMPBar.Drawable = common.Rectangle{}
	// leviMPBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	// leviMPBar.SetZIndex(4)
	// leviMPBar.Width = 83
	// leviMPBar.Height = 13
	// leviMPBar.Position = engo.Point{X: 410, Y: 271}
	// w.AddEntity(&leviMPBar)
	// //        MPBAR />
	// //        <CastBar
	// leviCBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviCBar.BasicEntity)
	// leviCBar.Drawable = common.Rectangle{}
	// leviCBar.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	// leviCBar.SetZIndex(4)
	// leviCBar.Width = 83
	// leviCBar.Height = 13
	// leviCBar.Position = engo.Point{X: 410, Y: 293}
	// w.AddEntity(&leviCBar)
	// //        CastBar />
	// //        <ACTION A
	// leviA := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviA.BasicEntity)
	// leviA.Drawable, _ = common.LoadedSprite(files[23])
	// leviA.Width = leviA.Drawable.Width()
	// leviA.Height = leviA.Drawable.Height()
	// leviA.SetZIndex(4)
	// leviA.SetCenter(engo.Point{X: 426, Y: 319})
	// w.AddEntity(&leviA)
	// //        ACTION A />
	// //        <ACTION B
	// leviB := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviB.BasicEntity)
	// leviB.Drawable, _ = common.LoadedSprite(files[21])
	// leviB.Width = leviB.Drawable.Width()
	// leviB.Height = leviB.Drawable.Height()
	// leviB.SetZIndex(4)
	// leviB.SetCenter(engo.Point{X: 477, Y: 319})
	// w.AddEntity(&leviB)
	// //        ACTION B />
	// //        <ACTION X
	// leviX := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviX.BasicEntity)
	// leviX.Drawable, _ = common.LoadedSprite(files[22])
	// leviX.Width = leviX.Drawable.Width()
	// leviX.Height = leviX.Drawable.Height()
	// leviX.SetZIndex(4)
	// leviX.SetCenter(engo.Point{X: 426, Y: 343})
	// w.AddEntity(&leviX)
	// //        ACTION X />
	// //        <ACTION Y
	// leviY := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// leviCard.AppendChild(&leviY.BasicEntity)
	// leviY.Drawable, _ = common.LoadedSprite(files[20])
	// leviY.Width = leviY.Drawable.Width()
	// leviY.Height = leviY.Drawable.Height()
	// leviY.SetZIndex(4)
	// leviY.SetCenter(engo.Point{X: 477, Y: 343})
	// w.AddEntity(&leviY)
	// //        ACTION Y />
	// // LEVI />
	//
	// // <WY
	// //    <Load font
	// wyFont := &common.Font{
	// 	Size: 64,
	// 	FG:   color.Black,
	// 	URL:  files[11],
	// }
	// wyFont.CreatePreloaded()
	// //    Load font />
	// //    < CARD
	// wyCard := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.Card = true
	// wyCard.Drawable, _ = common.LoadedSprite(files[7])
	// wyCard.SetZIndex(3)
	// wyCard.Position = engo.Point{X: 529, Y: 217}
	// w.AddEntity(&wyCard)
	// //    CARD />
	// //        <NAME
	// wyName := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyName.BasicEntity)
	// wyName.Drawable = common.Text{
	// 	Font: wyFont,
	// 	Text: "Wy",
	// }
	// wyName.SetZIndex(4)
	// wyName.Position = engo.Point{X: 535, Y: 222}
	// wyName.Scale = engo.Point{X: 0.25, Y: 0.25}
	// w.AddEntity(&wyName)
	// //        NAME />
	// //        <HPBAR
	// wyHPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyHPBar.AppendChild(&wyHPBar.BasicEntity)
	// wyHPBar.Drawable = common.Rectangle{}
	// wyHPBar.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	// wyHPBar.SetZIndex(4)
	// wyHPBar.Width = 83
	// wyHPBar.Height = 13
	// wyHPBar.Position = engo.Point{X: 537, Y: 249}
	// w.AddEntity(&wyHPBar)
	// //        HPBAR />
	// //        <MPBAR
	// wyMPBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyMPBar.BasicEntity)
	// wyMPBar.Drawable = common.Rectangle{}
	// wyMPBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	// wyMPBar.SetZIndex(4)
	// wyMPBar.Width = 83
	// wyMPBar.Height = 13
	// wyMPBar.Position = engo.Point{X: 537, Y: 271}
	// w.AddEntity(&wyMPBar)
	// //        MPBAR />
	// //        <CastBar
	// wyCBar := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyCBar.BasicEntity)
	// wyCBar.Drawable = common.Rectangle{}
	// wyCBar.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	// wyCBar.SetZIndex(4)
	// wyCBar.Width = 83
	// wyCBar.Height = 13
	// wyCBar.Position = engo.Point{X: 537, Y: 293}
	// w.AddEntity(&wyCBar)
	// //        CastBar />
	// //        <ACTION A
	// wyA := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyA.BasicEntity)
	// wyA.Drawable, _ = common.LoadedSprite(files[31])
	// wyA.Width = wyA.Drawable.Width()
	// wyA.Height = wyA.Drawable.Height()
	// wyA.SetZIndex(4)
	// wyA.SetCenter(engo.Point{X: 553, Y: 319})
	// w.AddEntity(&wyA)
	// //        ACTION A />
	// //        <ACTION B
	// wyB := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyB.BasicEntity)
	// wyB.Drawable, _ = common.LoadedSprite(files[29])
	// wyB.Width = wyB.Drawable.Width()
	// wyB.Height = wyB.Drawable.Height()
	// wyB.SetZIndex(4)
	// wyB.SetCenter(engo.Point{X: 604, Y: 319})
	// w.AddEntity(&wyB)
	// //        ACTION B />
	// //        <ACTION X
	// wyX := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyX.BasicEntity)
	// wyX.Drawable, _ = common.LoadedSprite(files[30])
	// wyX.Width = wyX.Drawable.Width()
	// wyX.Height = wyX.Drawable.Height()
	// wyX.SetZIndex(4)
	// wyX.SetCenter(engo.Point{X: 553, Y: 343})
	// w.AddEntity(&wyX)
	// //        ACTION X />
	// //        <ACTION Y
	// wyY := playerSelectableSprite{BasicEntity: ecs.NewBasic()}
	// wyCard.AppendChild(&wyY.BasicEntity)
	// wyY.Drawable, _ = common.LoadedSprite(files[28])
	// wyY.Width = wyY.Drawable.Width()
	// wyY.Height = wyY.Drawable.Height()
	// wyY.SetZIndex(4)
	// wyY.SetCenter(engo.Point{X: 604, Y: 343})
	// w.AddEntity(&wyY)
	// //        ACTION Y />
	// // WY />
	//
	// // < CursorComponentCallbacks
	//
	// // CursorComponentCallbacks />
}
