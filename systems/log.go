package systems

import (
	"sync"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type CombatLogComponent struct{}

func (c *CombatLogComponent) GetCombatLogComponent() *CombatLogComponent {
	return c
}

type CombatLogFace interface {
	GetCombatLogComponent() *CombatLogComponent
}

type CombatLogAble interface {
	CombatLogFace
	common.BasicFace
	common.RenderFace
}

type NotCombatLogComponent struct{}

func (n *NotCombatLogComponent) GetNotCombatLogComponent() *NotCombatLogComponent {
	return n
}

type NotCombatLogAble interface {
	GetNotCombatLogComponent() *NotCombatLogComponent
}

type CombatLogMessage struct {
	Msg  string
	Fnt  *common.Font
	Clip *common.Player
}

func (m CombatLogMessage) Type() string {
	return "CombatLogMessage"
}

type combatEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*CombatLogComponent
}

type CombatLogSystem struct {
	entities    []combatEntity
	lock        sync.RWMutex
	log         []CombatLogMessage
	idx, charAt int
	done, moved bool
	elapsed     float32
}

func (s *CombatLogSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("CombatLogMessage", func(message engo.Message) {
		msg, ok := message.(CombatLogMessage)
		if !ok {
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.log = append(s.log, msg)
	})
}

func (s *CombatLogSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, log *CombatLogComponent) {
	s.entities = append(s.entities, combatEntity{basic, render, log})
}

func (s *CombatLogSystem) AddByInterface(i ecs.Identifier) {
	a := i.(CombatLogAble)
	s.Add(a.GetBasicEntity(), a.GetRenderComponent(), a.GetCombatLogComponent())
}

func (s *CombatLogSystem) Remove(basic ecs.BasicEntity) {
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

func (s *CombatLogSystem) Update(dt float32) {
	s.elapsed += dt
	if s.done {
		if s.idx < len(s.log) {
			s.idx++
			s.moved = false
			s.done = false
		}
	} else {
		if !s.moved && len(s.log) > 0 {
			for i := len(s.entities) - 1; i > 0; i-- {
				s.entities[i].Drawable = s.entities[i-1].Drawable
			}
			txt := s.entities[0].Drawable.(common.Text)
			txt.Font = s.log[s.idx].Fnt
			txt.Text = ""
			s.entities[0].Drawable = txt
			s.moved = true
		}
		if len(s.log) > 0 && s.elapsed > 2 {
			s.log[s.idx].Clip.Play()
			s.charAt++
			txt := s.entities[0].Drawable.(common.Text)
			txt.Text = s.log[s.idx].Msg[:s.charAt]
			s.entities[0].Drawable = txt
			s.elapsed = 0
		}

		if len(s.log) > 0 && s.charAt >= len(s.log[s.idx].Msg) {
			s.charAt = 0
			s.elapsed = 0
			s.done = true
		}
	}
}
