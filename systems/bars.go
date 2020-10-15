package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type CharacterBarAble interface {
	common.BasicFace
	common.SpaceFace
	CharacterFace
	BarFace
}

type BaddieBarAble interface {
	common.BasicFace
	common.SpaceFace
	BaddieFace
	BarFace
}

type BarKind uint8

const (
	BarKindHP BarKind = iota
	BarKindMP
	BarKindCast
	BarKindBaddieCast
)

type BarComponent struct {
	Kind       BarKind
	TotalWidth float32
}

func (c *BarComponent) GetBarComponent() *BarComponent {
	return c
}

type BarFace interface {
	GetBarComponent() *BarComponent
}

type NotBarComponent struct{}

func (n *NotBarComponent) GetNotBarComponent() *NotBarComponent {
	return n
}

type NotBarAble interface {
	GetNotBarComponent() *NotBarComponent
}

type baddieBarEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*BarComponent
	*BaddieComponent
}

type characterBarEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*BarComponent
	*CharacterComponent
}

type BarSystem struct {
	bEntities []baddieBarEntity
	cEntities []characterBarEntity
}

func (s *BarSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, bar *BarComponent, baddie *BaddieComponent, chara *CharacterComponent) {
	if chara == nil {
		s.bEntities = append(s.bEntities, baddieBarEntity{basic, space, bar, baddie})
	} else {
		s.cEntities = append(s.cEntities, characterBarEntity{basic, space, bar, chara})
	}
}

func (s *BarSystem) AddByInterface(i ecs.Identifier) {
	b, ok := i.(BaddieBarAble)
	if !ok {
		c, ok := i.(CharacterBarAble)
		if !ok {
			return
		}
		s.Add(c.GetBasicEntity(), c.GetSpaceComponent(), c.GetBarComponent(), nil, c.GetCharacterComponent())
		return
	}
	s.Add(b.GetBasicEntity(), b.GetSpaceComponent(), b.GetBarComponent(), b.GetBaddieComponent(), nil)
}

func (s *BarSystem) Remove(b ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.bEntities {
		if entity.ID() == b.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.bEntities = append(s.bEntities[:delete], s.bEntities[delete+1:]...)
		return
	}
	for index, entity := range s.cEntities {
		if entity.ID() == b.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.cEntities = append(s.cEntities[:delete], s.cEntities[delete+1:]...)
	}
}

func (s *BarSystem) Update(dt float32) {
	for i := 0; i < len(s.cEntities); i++ {
		switch s.cEntities[i].Kind {
		case BarKindHP:
			s.cEntities[i].Width = s.cEntities[i].TotalWidth * float32(s.cEntities[i].HP) / float32(s.cEntities[i].MaxHP)
		case BarKindMP:
			s.cEntities[i].Width = s.cEntities[i].TotalWidth * float32(s.cEntities[i].MP) / float32(s.cEntities[i].MaxMP)
		case BarKindCast:
			s.cEntities[i].Width = s.cEntities[i].TotalWidth * float32(s.cEntities[i].CastTime) / float32(s.cEntities[i].CastAt)
		}
	}
	for i := 0; i < len(s.bEntities); i++ {
		if s.bEntities[i].CastTime != s.bEntities[i].CastAt {
			s.bEntities[i].Width = s.bEntities[i].TotalWidth * float32(s.bEntities[i].CastTime) / float32(s.bEntities[i].CastAt)
		} else {
			s.bEntities[i].Width = 0
		}
	}
}
