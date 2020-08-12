package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type PlayerSelectComponent struct {
	Selected, Card bool
}

func (c *PlayerSelectComponent) GetPlayerSelectComponent() *PlayerSelectComponent {
	return c
}

type PlayerSelectFace interface {
	GetPlayerSelectComponent() *PlayerSelectComponent
}

type PlayerSelectAble interface {
	PlayerSelectFace
	common.BasicFace
	common.SpaceFace
}

type NotPlayerSelectComponent struct{}

func (n *NotPlayerSelectComponent) GetNotPlayerSelectComponent() *NotPlayerSelectComponent {
	return n
}

type NotPlayerSelectAble interface {
	GetNotPlayerSelectComponent() *NotPlayerSelectComponent
}

type playerSelectEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*PlayerSelectComponent
}

type PlayerSelectSystem struct {
	adults   []playerSelectEntity
	children []playerSelectEntity
	idx, cur int
	Paused   bool
}

func (s *PlayerSelectSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, ps *PlayerSelectComponent) {
	if ps.Card {
		s.adults = append(s.adults, playerSelectEntity{basic, space, ps})
	} else {
		s.children = append(s.children, playerSelectEntity{basic, space, ps})
	}
}

func (s *PlayerSelectSystem) AddByInterface(i ecs.Identifier) {
	a := i.(PlayerSelectAble)
	s.Add(a.GetBasicEntity(), a.GetSpaceComponent(), a.GetPlayerSelectComponent())
}

func (s *PlayerSelectSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.adults {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.adults = append(s.adults[:delete], s.adults[delete+1:]...)
		return
	}
	delete = -1
	for index, entity := range s.children {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.children = append(s.children[:delete], s.children[delete+1:]...)
	}
}

func (s *PlayerSelectSystem) Update(dt float32) {
	if s.Paused {
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
	if s.idx >= len(s.adults) {
		s.idx = len(s.adults) - 1
	}
	if s.idx < 0 {
		s.idx = 0
	}
	if s.idx != s.cur {
		s.adults[s.cur].Position.Y += 25
		s.adults[s.cur].Selected = false
		for _, id := range s.adults[s.cur].Children() {
			for _, child := range s.children {
				if child.ID() == id.ID() {
					child.Position.Y += 25
					child.Selected = false
					break
				}
			}
		}
		s.adults[s.idx].Position.Y -= 25
		s.adults[s.idx].Selected = true
		for _, id := range s.adults[s.idx].Children() {
			for _, child := range s.children {
				if child.ID() == id.ID() {
					child.Position.Y -= 25
					child.Selected = true
					break
				}
			}
		}
		s.cur = s.idx
	}
}
