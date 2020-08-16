package scenes

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"

	"github.com/summersavers/summer/systems"
)

type sprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type selection struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	systems.CursorComponent
}

type audio struct {
	ecs.BasicEntity
	common.AudioComponent
}

type logText struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	systems.CombatLogComponent
}

type baddieSprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	systems.BaddieComponent
}

type playerSelectableSprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type selectionsceneswitch struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	systems.CursorComponent
	systems.SceneSwitchComponent
}

type character struct {
	ecs.BasicEntity
	systems.CharacterComponent
}
