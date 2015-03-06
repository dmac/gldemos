package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const WindowWidth = 800
const WindowHeight = 600

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Ocean Above", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	points := []float32{
		0, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}

	program, err := linkProgram("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	proj := mgl.Perspective(mgl.DegToRad(45.0), float32(WindowWidth)/WindowHeight, 0.1, 10.0)
	projUniform := gl.GetUniformLocation(program, gl.Str("proj\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &proj[0])

	camera := mgl.Vec3{0, 0, 5}
	view := mgl.LookAtV(camera, mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(1, 0, 0))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	vattrib := uint32(gl.GetAttribLocation(program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func linkProgram(vertexShaderFilename, fragmentShaderFilename string) (uint32, error) {
	vshader, err := compileShader("vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fshader, err := compileShader("fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vshader)
	gl.AttachShader(program, fshader)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.TRUE {
		return program, nil
	}
	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("link program: %s", log)
}

func compileShader(filename string, shaderType uint32) (uint32, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	source = append(source, 0)
	shader := gl.CreateShader(shaderType)
	csource := gl.Str(string(source))
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.TRUE {
		return shader, nil
	}
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("compile shader: %s%s", source, log)
}
