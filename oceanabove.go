package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

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
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	program, err := linkProgram("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	initBlockBase(program)

	blocks := []*Block{
		NewBlock(0, 0, 0),
		NewBlock(-2, 0, 0),
		NewBlock(2, 0, -2),
		NewBlock(0, 3, 0),
	}

	proj := mgl.Perspective(mgl.DegToRad(45.0), float32(WindowWidth)/WindowHeight, 0.1, 100.0)
	projUniform := gl.GetUniformLocation(program, gl.Str("proj\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &proj[0])

	camera := NewCamera(0, 1.5, 5, program)

	cx, cy := window.GetCursorPos()
	lastTime := time.Now().UnixNano()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currTime := time.Now().UnixNano()
		dt := float32(currTime-lastTime) / 1e9
		lastTime = currTime

		if window.GetKey(glfw.KeyW) == glfw.Press {
			camera.MoveForward(dt)
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			camera.MoveBackward(dt)
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			camera.MoveLeft(dt)
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			camera.MoveRight(dt)
		}
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		cx2, cy2 := window.GetCursorPos()
		dcx, dcy := float32(cx2-cx), float32(cy2-cy)
		if dcx != 0 {
			cx = cx2
			camera.Yaw -= dcx * 0.5
			camera.moved = true
		}
		if dcy != 0 {
			cy = cy2
			camera.Pitch -= dcy * 0.5
			camera.moved = true
			if camera.Pitch > 90 {
				camera.Pitch = 90
			}
			if camera.Pitch < -90 {
				camera.Pitch = -90
			}
		}

		camera.Update(dt)
		camera.Draw()
		for _, block := range blocks {
			block.Draw()
		}

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
