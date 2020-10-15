package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type MoveMessage struct {
	Char   *CharacterComponent
	Move   string
	Target *BaddieComponent
}

func (MoveMessage) Type() string {
	return "Move Message"
}

type MoveComponent struct {
	Elapsed   float32
	MovePhase int
	MoveFunc  func(s *MoveSystem, ent moveEntity, dt float32)
}

func (c *MoveComponent) GetMoveComponent() *MoveComponent {
	return c
}

type MoveFace interface {
	GetMoveComponent() *MoveComponent
}

type moveEntity struct {
	*ecs.BasicEntity
	*CharacterComponent
	*BaddieComponent
	*MoveComponent
}

var moves = map[string]func(s *MoveSystem, ent moveEntity, dt float32){
	"Regular Attack": func(s *MoveSystem, ent moveEntity, dt float32) {
		switch ent.MoveComponent.MovePhase {
		case 0:
			ent.CharacterComponent.CastTime = 0
			ent.CharacterComponent.CastAt = 3
			ent.CharacterComponent.MP -= 25
			ent.MoveComponent.MovePhase++
			ent.MoveComponent.Elapsed = 0
		case 1:
			ent.MoveComponent.Elapsed += dt
			if ent.MoveComponent.Elapsed >= 3 {
				message := ent.CharacterComponent.Name + " " + ent.CharacterComponent.AttackVerb + " the " + ent.BaddieComponent.Name + " for 15 damage!"
				engo.Mailbox.Dispatch(CombatLogMessage{
					Msg:  message,
					Fnt:  ent.CharacterComponent.Font,
					Clip: ent.CharacterComponent.Clip,
				})
				ent.BaddieComponent.HP -= 15
				ent.MoveComponent.MovePhase++
			}
		default:
			ent.CharacterComponent.CastTime = 0
			ent.CharacterComponent.CastAt = 0
			ent.MoveComponent.MovePhase = 0
			ent.MoveComponent.Elapsed = 0
			s.Remove(*ent.BasicEntity)
		}
	},
}

type MoveSystem struct {
	w        *ecs.World
	entities []moveEntity
}

func (s *MoveSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("Move Message", func(m engo.Message) {
		msg, ok := m.(MoveMessage)
		if !ok {
			return
		}
		basic := ecs.NewBasic()
		s.Add(&basic, msg.Char, msg.Target, &MoveComponent{MoveFunc: moves[msg.Move]})
	})
}

func (s *MoveSystem) Add(basic *ecs.BasicEntity, chara *CharacterComponent, baddie *BaddieComponent, move *MoveComponent) {
	s.entities = append(s.entities, moveEntity{basic, chara, baddie, move})
}

func (s *MoveSystem) Remove(basic ecs.BasicEntity) {
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

func (s *MoveSystem) Update(dt float32) {
	for _, e := range s.entities {
		if e.MoveComponent.MoveFunc != nil {
			e.MoveComponent.MoveFunc(s, e, dt)
		}
	}
}
