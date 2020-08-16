package systems

import "github.com/EngoEngine/engo/common"

type CharacterComponent struct {
	Name         string
	Abilities    []BattleBoxText
	Items        []BattleBoxText
	Acts         []BattleBoxText
	Font         *common.Font
	Clip         *common.Player
	Card         *common.Texture
	BattleBox    common.Drawable
	AIcon, BIcon *common.Texture
	XIcon, YIcon *common.Texture
	HP, MaxHP    int
	MP, MaxMP    int
	CastTime     int
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
