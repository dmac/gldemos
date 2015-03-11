package main

import (
	"math"

	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Pos   mgl.Vec3
	Vel   mgl.Vec3
	Acc   mgl.Vec3
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

func (c *Camera) MoveForward(dt float32) {
	c.Acc[0] -= float32(math.Sin(float64(mgl.DegToRad(c.Yaw))))
	c.Acc[2] -= float32(math.Cos(float64(mgl.DegToRad(c.Yaw))))
}

func (c *Camera) MoveBackward(dt float32) {
	c.Acc[0] += float32(math.Sin(float64(mgl.DegToRad(c.Yaw))))
	c.Acc[2] += float32(math.Cos(float64(mgl.DegToRad(c.Yaw))))
}

func (c *Camera) MoveLeft(dt float32) {
	c.Acc[0] -= float32(math.Sin(float64(mgl.DegToRad(c.Yaw + 90))))
	c.Acc[2] -= float32(math.Cos(float64(mgl.DegToRad(c.Yaw + 90))))
}

func (c *Camera) MoveRight(dt float32) {
	c.Acc[0] += float32(math.Sin(float64(mgl.DegToRad(c.Yaw + 90))))
	c.Acc[2] += float32(math.Cos(float64(mgl.DegToRad(c.Yaw + 90))))
}

func (c *Camera) Update(dt float32) {
	c.Pos = c.Pos.Add(c.Acc.Mul(dt * dt * 0.5).Add(c.Vel.Mul(dt)))

	drag := float32(10)
	if c.Acc.Len() > 0 {
		c.Acc = c.Acc.Normalize().Mul(drag * c.Speed)
	}

	//c.Acc[1] += -9.8 // gravity

	// https://stackoverflow.com/questions/667034/simple-physics-based-movement
	// v = a*dt - v0*drag*dt + v0
	c.Vel[0] = c.Vel[0] + c.Acc[0]*dt - c.Vel[0]*drag*dt
	c.Vel[1] = c.Vel[1] + c.Acc[1]*dt
	c.Vel[2] = c.Vel[2] + c.Acc[2]*dt - c.Vel[2]*drag*dt

	c.Acc = mgl.Vec3{0, 0, 0}
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
