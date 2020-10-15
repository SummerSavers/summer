package systems

import "github.com/EngoEngine/engo/common"

type CharacterComponent struct {
	Name         string
	AttackVerb   string
	Abilities    []BattleBoxText
	Items        []BattleBoxText
	Acts         []BattleBoxText
	Font         *common.Font
	Clip         *common.Player
	Card         *common.Texture
	BattleBox    common.Drawable
	AIcon, BIcon *common.Texture
	XIcon, YIcon *common.Texture
	HP, MaxHP    float32
	MP, MaxMP    float32
	CastAt       float32
	CastTime     float32
	CardSelected bool
}

type BattleBoxText struct {
	Name, Desc string
}

func (c *CharacterComponent) GetCharacterComponent() *CharacterComponent {
	return c
}

type CharacterFace interface {
	GetCharacterComponent() *CharacterComponent
}
