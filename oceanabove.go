package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

	camera := mgl.Vec3{0, 1.5, 5}
	camyaw := float32(0)
	campitch := float32(0)
	speed := float32(10)
	view := viewMatrix(camera, camyaw, campitch)
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(0, 0, 0))
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	cx, cy := window.GetCursorPos()
	lastTime := time.Now().UnixNano()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currTime := time.Now().UnixNano()
		dt := float32(currTime-lastTime) / 1e9
		lastTime = currTime

		moved := false
		if window.GetKey(glfw.KeyW) == glfw.Press {
			camera[0] -= dt * speed * float32(math.Sin(float64(mgl.DegToRad(camyaw))))
			camera[2] -= dt * speed * float32(math.Cos(float64(mgl.DegToRad(camyaw))))
			moved = true
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			camera[0] += dt * speed * float32(math.Sin(float64(mgl.DegToRad(camyaw))))
			camera[2] += dt * speed * float32(math.Cos(float64(mgl.DegToRad(camyaw))))
			moved = true
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			camera[0] -= dt * speed * float32(math.Sin(float64(mgl.DegToRad(camyaw+90))))
			camera[2] -= dt * speed * float32(math.Cos(float64(mgl.DegToRad(camyaw+90))))
			moved = true
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			camera[0] += dt * speed * float32(math.Sin(float64(mgl.DegToRad(camyaw+90))))
			camera[2] += dt * speed * float32(math.Cos(float64(mgl.DegToRad(camyaw+90))))
			moved = true
		}
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		cx2, cy2 := window.GetCursorPos()
		dcx, dcy := float32(cx2-cx), float32(cy2-cy)
		if dcx != 0 {
			cx = cx2
			camyaw -= dcx * 0.5
			moved = true
		}
		if dcy != 0 {
			cy = cy2
			campitch -= dcy * 0.5
			moved = true
			if campitch > 90 {
				campitch = 90
			}
			if campitch < -90 {
				campitch = -90
			}
		}

		if moved {
			view := viewMatrix(camera, camyaw, campitch)
			gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		}

		for _, block := range blocks {
			block.Draw()
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func viewMatrix(camera [3]float32, yaw float32, pitch float32) mgl.Mat4 {
	R := mgl.Rotate3DX(mgl.DegToRad(-pitch)).Mul3(mgl.Rotate3DY(mgl.DegToRad(-yaw))).Mat4()
	T := mgl.Translate3D(-camera[0], -camera[1], -camera[2])
	return R.Mul4(T)
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
