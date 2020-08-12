package shaders

import "github.com/EngoEngine/gl"

type ShaderDrawable struct{}

func (d ShaderDrawable) Texture() *gl.Texture {
	return nil
}

func (d ShaderDrawable) Width() float32 {
	return 0
}

func (d ShaderDrawable) Height() float32 {
	return 0
}

func (d ShaderDrawable) View() (float32, float32, float32, float32) {
	return 0, 0, 0, 0
}

func (d ShaderDrawable) Close() {}
