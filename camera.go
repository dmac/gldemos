package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Pos   mgl.Vec3
	Pitch float32
	Yaw   float32
	Speed float32

	moved       bool
	program     uint32
	viewUniform int32
}

func NewCamera(x, y, z float32, program uint32) *Camera {
	return &Camera{
		Pos:         mgl.Vec3{x, y, z},
		Pitch:       0,
		Yaw:         0,
		Speed:       10,
		moved:       true,
		program:     program,
		viewUniform: gl.GetUniformLocation(program, gl.Str("view\x00")),
	}
}

func (c *Camera) Draw() {
	if c.moved {
		view := c.viewMatrix()
		gl.UniformMatrix4fv(c.viewUniform, 1, false, &view[0])
	}
}

func (c *Camera) viewMatrix() mgl.Mat4 {
	R := mgl.Rotate3DX(mgl.DegToRad(-c.Pitch)).Mul3(mgl.Rotate3DY(mgl.DegToRad(-c.Yaw))).Mat4()
	T := mgl.Translate3D(-c.Pos[0], -c.Pos[1], -c.Pos[2])
	return R.Mul4(T)
}
