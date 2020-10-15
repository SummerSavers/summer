package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/math"
)

type TargetMessage struct {
	Character *CharacterComponent
	Move      string
	Target    TargetType
}

func (TargetMessage) Type() string {
	return "TargetMessage"
}

type TargetType uint8

const (
	TargetTypeBaddies TargetType = iota
	TargetTypeCharacters
	TargetTypeAllBaddies
	TargetTypeAllCharacters
	TargetTypeEverybody
)

type targetEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*BaddieComponent
}

type TargetSystem struct {
	entities      []targetEntity
	idx           int
	paused        bool
	elapsed       float32
	mover         *CharacterComponent
	move          string
	skipnextframe bool
}

func (s *TargetSystem) New(w *ecs.World) {
	s.paused = true
	engo.Mailbox.Listen("TargetMessage", func(m engo.Message) {
		msg, ok := m.(TargetMessage)
		if !ok || !s.paused { //only accept message when system is paused
			return
		}
		if msg.Target == TargetTypeBaddies {
			s.idx = 0
			s.paused = false
			s.skipnextframe = true
			s.mover = msg.Character
			s.move = msg.Move
		} else {
			engo.Mailbox.Dispatch(MoveMessage{
				Char:   msg.Character,
				Target: nil,
				Move:   msg.Move,
			})
		}
	})
}

func (s *TargetSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, baddie *BaddieComponent) {
	s.entities = append(s.entities, targetEntity{basic, render, baddie})
}

func (s *TargetSystem) Remove(basic ecs.BasicEntity) {
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

func (s *TargetSystem) Update(dt float32) {
	if s.skipnextframe {
		s.skipnextframe = false
		return
	}
	if s.paused {
		return
	}
	s.elapsed += dt
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
	} else if s.idx <= 0 {
		s.idx = 0
	}
	red := uint8(float32(0xa0)*math.Sin(s.elapsed*20)) + 0x55
	blue := uint8(float32(0xa0)*math.Cos(s.elapsed*20)) + 0x55
	s.entities[s.idx].Color = color.RGBA{R: red, G: 0x00, B: blue, A: 0xff}
	if engo.Input.Button("A").JustPressed() {
		engo.Mailbox.Dispatch(MoveMessage{
			Char:   s.mover,
			Target: s.entities[s.idx].BaddieComponent,
			Move:   s.move,
		})
		engo.Mailbox.Dispatch(PlayerSelectUnpauseMessage{})
		s.entities[s.idx].Color = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		s.paused = true
		s.mover = nil
		s.move = ""
	}
	if engo.Input.Button("B").JustPressed() {
		engo.Mailbox.Dispatch(PlayerSelectUnpauseMessage{})
		s.entities[s.idx].Color = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		s.paused = true
		s.mover = nil
		s.move = ""
	}

}
