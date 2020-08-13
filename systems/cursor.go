package systems

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type cursorChangeMessage struct{}

func (cursorChangeMessage) Type() string {
	return "Cursor Change Message"
}

type CursorEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
	*CursorComponent
}

type pointer struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type CursorComponent struct {
	ACallback, BCallback, XCallback, YCallback func(s *CursorSystem)
	// this is for lookups within the functions
	Classes []string
	enabled bool
}

func (c *CursorComponent) Enable() {
	c.enabled = true
}

func (c *CursorComponent) Disable() {
	c.enabled = false
	engo.Mailbox.Dispatch(cursorChangeMessage{})
}

func (c *CursorComponent) GetCursorComponent() *CursorComponent {
	return c
}

type NotCursorComponent struct{}

func (n *NotCursorComponent) GetNotCursorComponent() *NotCursorComponent {
	return n
}

type CursorSystem struct {
	entities []CursorEntity
	ptr      pointer
	indexAt  int
}

func (s *CursorSystem) New(w *ecs.World) {
	s.ptr = pointer{BasicEntity: ecs.NewBasic()}
	pointerTex, err := common.LoadedSprite("title/cursor.png")
	if err != nil {
		log.Fatalf("Unable to load pointer.png Error was: %v", err)
	}
	s.ptr.Drawable = pointerTex
	s.ptr.Hidden = true
	s.ptr.Width = s.ptr.Drawable.Width()
	s.ptr.Height = s.ptr.Drawable.Height()
	s.ptr.SetZIndex(100)
	w.AddEntity(&s.ptr)
	engo.Mailbox.Listen("Cursor Change Message", func(m engo.Message) {
		_, ok := m.(cursorChangeMessage)
		if !ok {
			return
		}
		if !s.entities[s.indexAt].enabled {
			var found bool
			for i := 0; i < len(s.entities); i++ {
				s.indexAt++
				if s.indexAt >= len(s.entities) {
					s.indexAt = 0
				}
				if s.entities[s.indexAt].enabled {
					found = true
					break
				}
			}
			if !found {
				s.ptr.Hidden = true
			}
		}
	})
}

func (s *CursorSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, render *common.RenderComponent, selection *CursorComponent) {
	if selection.ACallback == nil {
		selection.ACallback = func(*CursorSystem) {}
	}
	if selection.BCallback == nil {
		selection.BCallback = func(*CursorSystem) {}
	}
	if selection.XCallback == nil {
		selection.XCallback = func(*CursorSystem) {}
	}
	if selection.YCallback == nil {
		selection.YCallback = func(*CursorSystem) {}
	}
	selection.enabled = true
	s.entities = append(s.entities, CursorEntity{basic, space, render, selection})
	if len(s.entities) == 1 {
		s.SetPointer(0)
	}
}

func (s *CursorSystem) AddByInterface(id ecs.Identifier) {
	o := id.(CursorAble)
	s.Add(o.GetBasicEntity(), o.GetSpaceComponent(), o.GetRenderComponent(), o.GetCursorComponent())
}

func (s *CursorSystem) Remove(basic ecs.BasicEntity) {
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
	if s.indexAt >= len(s.entities) {
		s.indexAt = len(s.entities) - 1
		s.SetPointer(s.indexAt)
	}
}

func (s *CursorSystem) Update(dt float32) {
	var found bool
	for i := 0; i < len(s.entities); i++ {
		if engo.Input.Button("up").JustPressed() || engo.Input.Button("right").JustPressed() {
			s.indexAt++
			if s.indexAt >= len(s.entities) {
				s.indexAt = 0
			}
		} else if engo.Input.Button("down").JustPressed() || engo.Input.Button("left").JustPressed() {
			s.indexAt--
			if s.indexAt < 0 {
				s.indexAt = len(s.entities) - 1
			}
		}
		if s.entities[s.indexAt].enabled {
			found = true
			break
		}
	}
	if !found {
		return
	}
	s.SetPointer(s.indexAt)
	if engo.Input.Button("A").JustPressed() {
		s.entities[s.indexAt].ACallback(s)
	} else if engo.Input.Button("B").JustPressed() {
		s.entities[s.indexAt].BCallback(s)
	} else if engo.Input.Button("X").JustPressed() {
		s.entities[s.indexAt].XCallback(s)
	} else if engo.Input.Button("Y").JustPressed() {
		s.entities[s.indexAt].YCallback(s)
	}
}

func (s *CursorSystem) SetPointer(i int) {
	if len(s.entities) == 0 || i > len(s.entities)-1 {
		s.ptr.Hidden = true
		return
	}
	ent := s.entities[i]
	s.ptr.Hidden = false
	s.ptr.Position.X = ent.Position.X - s.ptr.Width - 2
	s.ptr.Position.Y = ent.Position.Y + (ent.Height / 2) + 6
}

type CursorFace interface {
	GetCursorComponent() *CursorComponent
}

type CursorAble interface {
	common.BasicFace
	common.SpaceFace
	common.RenderFace
	CursorFace
}

type NotCursorAble interface {
	GetNotCursorComponent() *NotCursorComponent
}
