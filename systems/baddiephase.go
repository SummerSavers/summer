package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type BaddieComponent struct {
	Name         string
	Phases       map[string]Phase
	CurrentPhase string
	Font         *common.Font
	Clip         *common.Player
	HP, MaxHP    int
	CastAt       int
	CastTime     int
	Spritesheet  *common.Spritesheet
}

type Phase struct {
	Animations []*common.Animation
	Abilities  []BaddieAbility
}

type BaddieAbility struct {
	Weight      int
	LogMessages []string
	Times       []float32
}

func (c *BaddieComponent) GetBaddieComponent() *BaddieComponent {
	return c
}

type BaddieFace interface {
	GetBaddieComponent() *BaddieComponent
}

type BaddiePhaseAble interface {
	BaddieFace
	common.BasicFace
	common.RenderFace
}

type NotBaddiePhaseComponent struct{}

func (n *NotBaddiePhaseComponent) GetNotBaddiePhaseComponent() *NotBaddiePhaseComponent {
	return n
}

type NotBaddiePhaseAble interface {
	GetNotBaddiePhaseComponent() *NotBaddiePhaseComponent
}

type BaddiePhaseMessage struct {
	ID        uint64
	NewSprite common.Drawable
}

func (m BaddiePhaseMessage) Type() string {
	return "BaddiePhaseMessage"
}

type baddiePhaseEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*BaddieComponent
}

type BaddiePhaseSystem struct {
	entities []baddiePhaseEntity
}

func (s *BaddiePhaseSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("BaddiePhaseMessage", func(message engo.Message) {
		msg, ok := message.(BaddiePhaseMessage)
		if !ok {
			return
		}
		for _, e := range s.entities {
			if e.ID() == msg.ID {
				e.Drawable = msg.NewSprite
			}
		}
	})
}

func (s *BaddiePhaseSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, bp *BaddieComponent) {
	s.entities = append(s.entities, baddiePhaseEntity{basic, render, bp})
}

func (s *BaddiePhaseSystem) AddByInterface(i ecs.Identifier) {
	a := i.(BaddiePhaseAble)
	s.Add(a.GetBasicEntity(), a.GetRenderComponent(), a.GetBaddieComponent())
}

func (s *BaddiePhaseSystem) Remove(basic ecs.BasicEntity) {
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

func (s *BaddiePhaseSystem) Update(dt float32) {}
