package systems

import (
	"bytes"
	"image"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"golang.org/x/image/font/gofont/gomono"
)

type BattleboxAble interface {
	common.BasicFace
	CharacterFace
}

type NotBattleboxComponent struct{}

func (c *NotBattleboxComponent) GetNotBattleboComponent() *NotBattleboxComponent {
	return c
}

type NotBattleboxAble interface {
	GetNotBattleboxComponent() *NotBattleboxComponent
}

type BattleBoxMenu uint8

const (
	BattleBoxMenuAbilities = iota
	BattleBoxMenuItems
	BattleBoxMenuActs
)

type BattleBoxShowMessage struct {
	MenuToShow BattleBoxMenu
}

func (BattleBoxShowMessage) Type() string {
	return "battleboxshowmessage"
}

type BattleBoxHideMessage struct{}

func (BattleBoxHideMessage) Type() string {
	return "battleboxhidemessage"
}

type battleboxEntity struct {
	*ecs.BasicEntity
	*CharacterComponent
}

type BattleboxSystem struct {
	entities                   []battleboxEntity
	battlebox                  sprite
	item1, item2, item3, item4 selection
	desc                       sprite
	cur                        *CursorSystem
	idx                        int
	paused                     bool
	battleboxMenuSelected      BattleBoxMenu
}

func (s *BattleboxSystem) New(w *ecs.World) {
	engo.Files.LoadReaderData("gomono.ttf", bytes.NewReader(gomono.TTF))
	fnt := &common.Font{
		URL:  "gomono.ttf",
		FG:   color.Black,
		Size: 64,
	}
	fnt.CreatePreloaded()
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.cur = sys
		}
	}

	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0xff, 0x00, 0x00, 0xff})
	imgobj := common.NewImageObject(img)
	tex := common.NewTextureSingle(imgobj)

	//BattleBox
	s.battlebox = sprite{BasicEntity: ecs.NewBasic()}
	s.battlebox.Drawable = tex
	s.battlebox.SetZIndex(5)
	s.battlebox.Hidden = true
	s.battlebox.Position.Set(20, 150)
	s.battlebox.Color = color.RGBA{0xff, 0xff, 0xff, 0xe5}
	w.AddEntity(&s.battlebox)
	//item1
	s.item1 = selection{BasicEntity: ecs.NewBasic()}
	s.item1.Drawable = common.Text{
		Font: fnt,
		Text: "---",
	}
	s.item1.SetZIndex(6)
	s.item1.Scale = engo.Point{X: 0.4, Y: 0.4}
	s.item1.Position = engo.Point{X: 55, Y: 175}
	s.item1.Hidden = true
	w.AddEntity(&s.item1)
	//Item2
	s.item2 = selection{BasicEntity: ecs.NewBasic()}
	s.item2.Drawable = common.Text{
		Font: fnt,
		Text: "---",
	}
	s.item2.SetZIndex(6)
	s.item2.Scale = engo.Point{X: 0.4, Y: 0.4}
	s.item2.Position = engo.Point{X: 355, Y: 175}
	s.item2.Hidden = true
	w.AddEntity(&s.item2)
	//Item3
	s.item3 = selection{BasicEntity: ecs.NewBasic()}
	s.item3.Drawable = common.Text{
		Font: fnt,
		Text: "---",
	}
	s.item3.SetZIndex(6)
	s.item3.Scale = engo.Point{X: 0.4, Y: 0.4}
	s.item3.Position = engo.Point{X: 55, Y: 235}
	s.item3.Hidden = true
	w.AddEntity(&s.item3)
	//Item4
	s.item4 = selection{BasicEntity: ecs.NewBasic()}
	s.item4.Drawable = common.Text{
		Font: fnt,
		Text: "---",
	}
	s.item4.SetZIndex(6)
	s.item4.Scale = engo.Point{X: 0.4, Y: 0.4}
	s.item4.Position = engo.Point{X: 355, Y: 235}
	s.item4.Hidden = true
	w.AddEntity(&s.item4)
	//DescLine1
	s.desc = sprite{BasicEntity: ecs.NewBasic()}
	s.desc.Drawable = common.Text{
		Font: fnt,
		Text: "---",
	}
	s.desc.SetZIndex(6)
	s.desc.Scale = engo.Point{X: 0.4, Y: 0.4}
	s.desc.Position = engo.Point{X: 355, Y: 235}
	s.desc.Hidden = true
	w.AddEntity(&s.desc)
	s.paused = true
	engo.Mailbox.Listen("battleboxhidemessage", func(msg engo.Message) {
		_, ok := msg.(BattleBoxHideMessage)
		if !ok {
			return
		}
		s.battlebox.Hidden = true
		s.desc.Hidden = true
		s.item1.Hidden = true
		s.item2.Hidden = true
		s.item3.Hidden = true
		s.item4.Hidden = true
		s.item1.Selected = false
		s.item2.Selected = false
		s.item3.Selected = false
		s.item4.Selected = false
		s.cur.Remove(s.item1.BasicEntity)
		s.cur.Remove(s.item2.BasicEntity)
		s.cur.Remove(s.item3.BasicEntity)
		s.cur.Remove(s.item4.BasicEntity)
		engo.Mailbox.Dispatch(CursorSetMessage{-1})
		s.idx = 0
		s.paused = true
	})
	engo.Mailbox.Listen("battleboxshowmessage", func(msg engo.Message) {
		m, ok := msg.(BattleBoxShowMessage)
		if !ok {
			return
		}
		s.battlebox.Hidden = false
		s.desc.Hidden = false
		s.item1.Hidden = false
		s.item2.Hidden = false
		s.item3.Hidden = false
		s.item4.Hidden = false
		s.item1.Selected = true
		s.cur.AddByInterface(s.item1)
		s.cur.AddByInterface(s.item2)
		s.cur.AddByInterface(s.item3)
		s.cur.AddByInterface(s.item4)
		s.idx = 0
		for i := 0; i < len(s.entities); i++ {
			if s.entities[i].CardSelected {
				s.battlebox.Drawable = s.entities[i].BattleBox
				desc := "..."
				switch m.MenuToShow {
				case BattleBoxMenuAbilities:
					s.battleboxMenuSelected = BattleBoxMenuAbilities
					if len(s.entities[i].Abilities) > 0 {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Abilities[0].Name,
						}
					} else {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Abilities) > 1 {
						desc = s.entities[i].Abilities[0].Desc
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Abilities[1].Name,
						}
					} else {
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Abilities) > 2 {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Abilities[2].Name,
						}
					} else {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Abilities) > 3 {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Abilities[3].Name,
						}
					} else {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					s.desc.Drawable = common.Text{
						Font: s.entities[i].Font,
						Text: desc,
					}
				case BattleBoxMenuItems:
					s.battleboxMenuSelected = BattleBoxMenuItems
					if len(s.entities[i].Items) > 0 {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Items[0].Name,
						}
					} else {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Items) > 1 {
						desc = s.entities[i].Items[0].Desc
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Items[1].Name,
						}
					} else {
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Items) > 2 {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Items[2].Name,
						}
					} else {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Items) > 3 {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Items[3].Name,
						}
					} else {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					s.desc.Drawable = common.Text{
						Font: s.entities[i].Font,
						Text: desc,
					}
				case BattleBoxMenuActs:
					s.battleboxMenuSelected = BattleBoxMenuItems
					if len(s.entities[i].Acts) > 0 {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Acts[0].Name,
						}
					} else {
						s.item1.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Acts) > 1 {
						desc = s.entities[i].Acts[0].Desc
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Acts[1].Name,
						}
					} else {
						s.item2.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Acts) > 2 {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Acts[2].Name,
						}
					} else {
						s.item3.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					if len(s.entities[i].Acts) > 3 {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: s.entities[i].Acts[3].Name,
						}
					} else {
						s.item4.Drawable = common.Text{
							Font: s.entities[i].Font,
							Text: "---",
						}
					}
					s.desc.Drawable = common.Text{
						Font: s.entities[i].Font,
						Text: desc,
					}
				}
			}
		}
		s.paused = false
	})
}

func (s *BattleboxSystem) Add(basic *ecs.BasicEntity, character *CharacterComponent) {
	s.entities = append(s.entities, battleboxEntity{basic, character})
}

func (s *BattleboxSystem) AddByInterface(id ecs.Identifier) {
	o, ok := id.(BattleboxAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetCharacterComponent())
}

func (s *BattleboxSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *BattleboxSystem) Update(dt float32) {
	if s.paused {
		return
	}
	cur := s.idx

	entIdx := -1
	for i := 0; i < len(s.entities); i++ {
		if s.entities[i].CardSelected {
			entIdx = i
		}
	}
	if entIdx < 0 {
		return
	}

	if engo.Input.Button("down").JustPressed() || engo.Input.Button("left").JustPressed() {
		s.idx++
		l := 0
		switch s.battleboxMenuSelected {
		case BattleBoxMenuAbilities:
			l = len(s.entities[entIdx].Abilities)
		case BattleBoxMenuItems:
			l = len(s.entities[entIdx].Items)
		case BattleBoxMenuActs:
			l = len(s.entities[entIdx].Acts)
		}
		if s.idx > l-1 {
			s.idx = l - 1
		}
	} else if engo.Input.Button("up").JustPressed() || engo.Input.Button("right").JustPressed() {
		s.idx--
		if s.idx < 0 {
			s.idx = 0
		}
	}

	if cur != s.idx {
		if s.idx < cur && s.item1.Selected {
			txt1 := s.item1.Drawable.(common.Text)
			txt2 := s.item2.Drawable.(common.Text)
			txt3 := s.item3.Drawable.(common.Text)
			txt := ""
			switch s.battleboxMenuSelected {
			case BattleBoxMenuAbilities:
				txt = s.entities[entIdx].Abilities[s.idx].Name
			case BattleBoxMenuItems:
				txt = s.entities[entIdx].Items[s.idx].Name
			case BattleBoxMenuActs:
				txt = s.entities[entIdx].Acts[s.idx].Name
			}
			s.item1.Drawable = common.Text{
				Font: txt1.Font,
				Text: txt,
			}
			s.item2.Drawable = txt1
			s.item3.Drawable = txt2
			s.item4.Drawable = txt3
		} else if s.idx > cur && s.item4.Selected {
			txt2 := s.item2.Drawable.(common.Text)
			txt3 := s.item3.Drawable.(common.Text)
			txt4 := s.item4.Drawable.(common.Text)
			txt := ""
			switch s.battleboxMenuSelected {
			case BattleBoxMenuAbilities:
				txt = s.entities[entIdx].Abilities[s.idx].Name
			case BattleBoxMenuItems:
				txt = s.entities[entIdx].Items[s.idx].Name
			case BattleBoxMenuActs:
				txt = s.entities[entIdx].Acts[s.idx].Name
			}
			s.item4.Drawable = common.Text{
				Font: txt4.Font,
				Text: txt,
			}
			s.item1.Drawable = txt2
			s.item2.Drawable = txt3
			s.item3.Drawable = txt4
		}
		txt := ""
		switch s.battleboxMenuSelected {
		case BattleBoxMenuAbilities:
			txt = s.entities[entIdx].Abilities[s.idx].Desc
		case BattleBoxMenuItems:
			txt = s.entities[entIdx].Items[s.idx].Desc
		case BattleBoxMenuActs:
			txt = s.entities[entIdx].Acts[s.idx].Desc
		}
		s.desc.Drawable = common.Text{
			Font: s.entities[entIdx].Font,
			Text: txt,
		}
	}

	if engo.Input.Button("A").JustPressed() {
		//message to target system which goes to attack system
		//depends on the state, luckily we know that ^_^
	}

	if engo.Input.Button("B").JustPressed() {
		//message to go back (close this and reactivate player select)
	}
}
