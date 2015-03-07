package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Block struct {
	Pos  mgl.Vec3
	Size float32

	// TODO(dmac) There should likely only be a single, shared instance of the mesh that's transformed
	// with a model matrix..
	vao uint32
}

func NewBlock(x, y, z float32, program uint32) *Block {
	block := &Block{
		Pos:  mgl.Vec3{x, y, z},
		Size: 1,
	}

	gl.GenVertexArrays(1, &block.vao)
	gl.BindVertexArray(block.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	mesh := block.genVertexData()
	gl.BufferData(gl.ARRAY_BUFFER, len(mesh)*4, gl.Ptr(mesh), gl.STATIC_DRAW)

	vattrib := uint32(gl.GetAttribLocation(program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	return block
}

// genVertexData takes an array of points and returns a "mesh" of the block for use in glBufferData.
func (b *Block) genVertexData() []float32 {
	ps := b.genPoints()
	return []float32{
		// top face
		ps[1][0], ps[1][1], ps[1][2],
		ps[2][0], ps[2][1], ps[2][2],
		ps[3][0], ps[3][1], ps[3][2],
		ps[3][0], ps[3][1], ps[3][2],
		ps[0][0], ps[0][1], ps[0][2],
		ps[1][0], ps[1][1], ps[1][2],

		// bottom face
		ps[4][0], ps[4][1], ps[4][2],
		ps[7][0], ps[7][1], ps[7][2],
		ps[6][0], ps[6][1], ps[6][2],
		ps[6][0], ps[6][1], ps[6][2],
		ps[5][0], ps[5][1], ps[5][2],
		ps[4][0], ps[4][1], ps[4][2],

		// front face
		ps[0][0], ps[0][1], ps[0][2],
		ps[3][0], ps[3][1], ps[3][2],
		ps[7][0], ps[7][1], ps[7][2],
		ps[7][0], ps[7][1], ps[7][2],
		ps[4][0], ps[4][1], ps[4][2],
		ps[0][0], ps[0][1], ps[0][2],

		// back face
		ps[1][0], ps[1][1], ps[1][2],
		ps[5][0], ps[5][1], ps[5][2],
		ps[6][0], ps[6][1], ps[6][2],
		ps[6][0], ps[6][1], ps[6][2],
		ps[2][0], ps[2][1], ps[2][2],
		ps[1][0], ps[1][1], ps[1][2],

		// left face
		ps[2][0], ps[2][1], ps[2][2],
		ps[6][0], ps[6][1], ps[6][2],
		ps[7][0], ps[7][1], ps[7][2],
		ps[7][0], ps[7][1], ps[7][2],
		ps[3][0], ps[3][1], ps[3][2],
		ps[2][0], ps[2][1], ps[2][2],

		// right face
		ps[0][0], ps[0][1], ps[0][2],
		ps[4][0], ps[4][1], ps[4][2],
		ps[5][0], ps[5][1], ps[5][2],
		ps[5][0], ps[5][1], ps[5][2],
		ps[1][0], ps[1][1], ps[1][2],
		ps[0][0], ps[0][1], ps[0][2],
	}
}

// genPoints calculates the eight corners of b using Pos and Size.
func (b *Block) genPoints() [8]mgl.Vec3 {
	halfsize := b.Size / 2
	center := mgl.Vec3{b.Pos[0], b.Pos[1], b.Pos[2]}
	v1 := center.Add(mgl.Vec3{halfsize, halfsize, halfsize})
	v2 := center.Add(mgl.Vec3{halfsize, halfsize, -halfsize})
	v3 := center.Add(mgl.Vec3{-halfsize, halfsize, -halfsize})
	v4 := center.Add(mgl.Vec3{-halfsize, halfsize, halfsize})
	v5 := center.Add(mgl.Vec3{halfsize, -halfsize, halfsize})
	v6 := center.Add(mgl.Vec3{halfsize, -halfsize, -halfsize})
	v7 := center.Add(mgl.Vec3{-halfsize, -halfsize, -halfsize})
	v8 := center.Add(mgl.Vec3{-halfsize, -halfsize, halfsize})
	return [8]mgl.Vec3{v1, v2, v3, v4, v5, v6, v7, v8}
}
