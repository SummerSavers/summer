package main

import (
	"github.com/summersavers/summer/scenes"

	"github.com/EngoEngine/engo"
)

func main() {
	title := scenes.TitleScene{}
	engo.RegisterScene(&title)
	engo.RegisterScene(&scenes.IntroScene{})
	engo.Run(engo.RunOptions{
		Title:                      "Let's Save Summer!",
		Width:                      640, //512, //16
		Height:                     360, //288, //9
		ScaleOnResize:              true,
		FPSLimit:                   60,
		ApplicationMajorVersion:    0,
		ApplicationMinorVersion:    1,
		ApplicationRevisionVersion: 0,
	}, &title)
}
