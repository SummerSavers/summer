package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type PlayerSelectAble interface {
	common.BasicFace
	CharacterFace
}

type NotPlayerSelectComponent struct{}

func (n *NotPlayerSelectComponent) GetNotPlayerSelectComponent() *NotPlayerSelectComponent {
	return n
}

type NotPlayerSelectAble interface {
	GetNotPlayerSelectComponent() *NotPlayerSelectComponent
}

type PlayerSelectUnpauseMessage struct{}

func (PlayerSelectUnpauseMessage) Type() string {
	return "playerselectunpausemessage"
}

type playerSelectEntity struct {
	*ecs.BasicEntity
	*CharacterComponent
}

type PlayerSelectSystem struct {
	PlayerCount   int
	cardpositions []engo.Point
	entities      []playerSelectEntity
	sprites       map[uint64][]*sprite
	w             *ecs.World
	idx, cur      int
	paused        bool
	skipNextFrame bool
}

func (s *PlayerSelectSystem) Priority() int {
	return 0
}

func (s *PlayerSelectSystem) New(w *ecs.World) {
	s.sprites = make(map[uint64][]*sprite)
	s.w = w
	s.setupPositions()

	engo.Mailbox.Listen("playerselectunpausemessage", func(msg engo.Message) {
		s.paused = false
		s.skipNextFrame = true
	})
}

func (s *PlayerSelectSystem) setupPositions() {
	switch s.PlayerCount {
	case 0:
		s.paused = true
	case 1:
		s.cardpositions = []engo.Point{
			engo.Point{X: 279, Y: 191},
		}
	case 2:
		s.cardpositions = []engo.Point{
			engo.Point{X: 172, Y: 191},
			engo.Point{X: 284, Y: 216},
		}
	case 3:
		s.cardpositions = []engo.Point{
			engo.Point{X: 157, Y: 191},
			engo.Point{X: 269, Y: 216},
			engo.Point{X: 381, Y: 216},
		}
	case 4:
		s.cardpositions = []engo.Point{
			engo.Point{X: 101, Y: 191},
			engo.Point{X: 213, Y: 216},
			engo.Point{X: 325, Y: 216},
			engo.Point{X: 437, Y: 216},
		}
	case 5:
		s.cardpositions = []engo.Point{
			engo.Point{X: 45, Y: 191},
			engo.Point{X: 157, Y: 216},
			engo.Point{X: 269, Y: 216},
			engo.Point{X: 381, Y: 216},
			engo.Point{X: 493, Y: 216},
		}
	}
}

func (s *PlayerSelectSystem) Add(basic *ecs.BasicEntity, chara *CharacterComponent) {
	sprs := make([]*sprite, 0)
	//create player card
	card := sprite{BasicEntity: ecs.NewBasic()}
	card.Drawable = chara.Card
	card.SetZIndex(3)
	card.Position = s.cardpositions[len(s.entities)]
	sprs = append(sprs, &card)
	s.w.AddEntity(&card)
	//create player  name
	name := sprite{BasicEntity: ecs.NewBasic()}
	name.Drawable = common.Text{
		Font: chara.Font,
		Text: chara.Name,
	}
	name.SetZIndex(4)
	name.Scale = engo.Point{X: 0.2, Y: 0.2}
	name.Position = s.cardpositions[len(s.entities)]
	name.Position.X += 6
	name.Position.Y += 7
	sprs = append(sprs, &name)
	s.w.AddEntity(&name)
	//create attack icon
	attack := sprite{BasicEntity: ecs.NewBasic()}
	attack.Drawable = chara.AIcon
	attack.SetZIndex(4)
	attack.Position = s.cardpositions[len(s.entities)]
	attack.Position.X += 13
	attack.Position.Y += 91
	sprs = append(sprs, &attack)
	s.w.AddEntity(&attack)
	//create ability icon
	ability := sprite{BasicEntity: ecs.NewBasic()}
	ability.Drawable = chara.BIcon
	ability.SetZIndex(4)
	ability.Position = s.cardpositions[len(s.entities)]
	ability.Position.X += 64
	ability.Position.Y += 91
	sprs = append(sprs, &ability)
	s.w.AddEntity(&ability)
	//create items icon
	item := sprite{BasicEntity: ecs.NewBasic()}
	item.Drawable = chara.XIcon
	item.SetZIndex(4)
	item.Position = s.cardpositions[len(s.entities)]
	item.Position.X += 13
	item.Position.Y += 115
	sprs = append(sprs, &item)
	s.w.AddEntity(&item)
	//create act icon
	act := sprite{BasicEntity: ecs.NewBasic()}
	act.Drawable = chara.YIcon
	act.SetZIndex(4)
	act.Position = s.cardpositions[len(s.entities)]
	act.Position.X += 64
	act.Position.Y += 115
	sprs = append(sprs, &act)
	s.w.AddEntity(&act)

	s.entities = append(s.entities, playerSelectEntity{basic, chara})
	s.sprites[basic.ID()] = sprs
}

func (s *PlayerSelectSystem) AddByInterface(i ecs.Identifier) {
	a := i.(PlayerSelectAble)
	s.Add(a.GetBasicEntity(), a.GetCharacterComponent())
}

func (s *PlayerSelectSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
		return
	}
}

func (s *PlayerSelectSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	if engo.Input.Button("up").JustPressed() {
		s.idx++
	}
	if engo.Input.Button("right").JustPressed() {
		s.idx++
	}
	if engo.Input.Button("down").JustPressed() {
		s.idx--
	}
	if engo.Input.Button("left").JustPressed() {
		s.idx--
	}
	if s.idx >= len(s.entities) {
		s.idx = len(s.entities) - 1
	}
	if s.idx < 0 {
		s.idx = 0
	}
	if s.idx != s.cur {
		for _, spr := range s.sprites[s.entities[s.cur].ID()] {
			spr.Position.Y += 25
		}
		s.entities[s.cur].CardSelected = false
		for _, spr := range s.sprites[s.entities[s.idx].ID()] {
			spr.Position.Y -= 25
		}
		s.entities[s.idx].CardSelected = true
		s.cur = s.idx
	}
	if engo.Input.Button("A").JustPressed() {
		//message to target system which goes to attack system
	}
	if engo.Input.Button("B").JustPressed() {
		engo.Mailbox.Dispatch(BattleBoxShowMessage{MenuToShow: BattleBoxMenuAbilities})
		s.skipNextFrame = true
		s.paused = true
	}
	if engo.Input.Button("X").JustPressed() {
		engo.Mailbox.Dispatch(BattleBoxShowMessage{MenuToShow: BattleBoxMenuItems})
		s.skipNextFrame = true
		s.paused = true
	}
	if engo.Input.Button("Y").JustPressed() {
		engo.Mailbox.Dispatch(BattleBoxShowMessage{MenuToShow: BattleBoxMenuActs})
		s.skipNextFrame = true
		s.paused = true
	}
}
