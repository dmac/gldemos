package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

var blockSpec BlockSpec

func initBlockSpec(program uint32) {
	blockSpec = BlockSpec{0, program}
	gl.GenVertexArrays(1, &blockSpec.vao)
	gl.BindVertexArray(blockSpec.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	mesh := genBlockMesh()
	gl.BufferData(gl.ARRAY_BUFFER, len(mesh)*4, gl.Ptr(mesh), gl.STATIC_DRAW)

	vattrib := uint32(gl.GetAttribLocation(program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
}

type BlockSpec struct {
	vao     uint32
	program uint32
}

type Block struct {
	BlockSpec
	Pos   mgl.Vec3
	Size  float32
	Pitch float32
	Yaw   float32
}

func NewBlock(x, y, z float32) *Block {
	block := &Block{
		BlockSpec: blockSpec,
		Pos:       mgl.Vec3{x, y, z},
		Size:      1,
		Pitch:     45,
		Yaw:       45,
	}
	return block
}

// TODO(dmac) Should only update the uniform when necessary (i.e., after movement)
func (b *Block) Draw() {
	model := b.modelMatrix()
	modelUniform := gl.GetUniformLocation(b.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	gl.BindVertexArray(b.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
}

// modelMatrix generates a model matrix to be used as a uniform in the shader program.
func (b *Block) modelMatrix() mgl.Mat4 {
	S := mgl.Ident3().Mul(b.Size).Mat4()
	T := mgl.Translate3D(b.Pos[0], b.Pos[1], b.Pos[2]).Mul(b.Size)
	R := mgl.Rotate3DX(mgl.DegToRad(b.Pitch)).Mul3(mgl.Rotate3DY(mgl.DegToRad(b.Yaw))).Mat4()
	return T.Mul4(R).Mul4(S)
}

// genBlockMesh takes an array of points and returns a "mesh" of the block for use in glBufferData.
func genBlockMesh() []float32 {
	p0 := mgl.Vec3{0.5, 0.5, 0.5}
	p1 := mgl.Vec3{0.5, 0.5, -0.5}
	p2 := mgl.Vec3{-0.5, 0.5, -0.5}
	p3 := mgl.Vec3{-0.5, 0.5, 0.5}
	p4 := mgl.Vec3{0.5, -0.5, 0.5}
	p5 := mgl.Vec3{0.5, -0.5, -0.5}
	p6 := mgl.Vec3{-0.5, -0.5, -0.5}
	p7 := mgl.Vec3{-0.5, -0.5, 0.5}
	return []float32{
		// top face
		p1[0], p1[1], p1[2],
		p2[0], p2[1], p2[2],
		p3[0], p3[1], p3[2],
		p3[0], p3[1], p3[2],
		p0[0], p0[1], p0[2],
		p1[0], p1[1], p1[2],

		// bottom face
		p4[0], p4[1], p4[2],
		p7[0], p7[1], p7[2],
		p6[0], p6[1], p6[2],
		p6[0], p6[1], p6[2],
		p5[0], p5[1], p5[2],
		p4[0], p4[1], p4[2],

		// front face
		p0[0], p0[1], p0[2],
		p3[0], p3[1], p3[2],
		p7[0], p7[1], p7[2],
		p7[0], p7[1], p7[2],
		p4[0], p4[1], p4[2],
		p0[0], p0[1], p0[2],

		// back face
		p1[0], p1[1], p1[2],
		p5[0], p5[1], p5[2],
		p6[0], p6[1], p6[2],
		p6[0], p6[1], p6[2],
		p2[0], p2[1], p2[2],
		p1[0], p1[1], p1[2],

		// left face
		p2[0], p2[1], p2[2],
		p6[0], p6[1], p6[2],
		p7[0], p7[1], p7[2],
		p7[0], p7[1], p7[2],
		p3[0], p3[1], p3[2],
		p2[0], p2[1], p2[2],

		// right face
		p0[0], p0[1], p0[2],
		p4[0], p4[1], p4[2],
		p5[0], p5[1], p5[2],
		p5[0], p5[1], p5[2],
		p1[0], p1[1], p1[2],
		p0[0], p0[1], p0[2],
	}
}
