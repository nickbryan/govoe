package renderer

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	vertexShader = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 clr; 

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1);
	clr = aColor;
}
` + "\x00"
	fragmentShader = `
#version 410 core
out vec4 FragColor;
in vec3 clr;

void main()
{
    FragColor = vec4(clr, 1.0f);
} 
` + "\x00"
)

type Renderer struct {
	vertexShader, fragmentShader uint32
	shaderProgram                uint32
	meshes                       []*Mesh
}

func New() *Renderer {
	r := &Renderer{}

	r.Setup() // TODO: move this to a once callback somewhere to ensure initialisation (update loop maybe?)

	return r
}

func (r *Renderer) Setup() {
	var nrAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	fmt.Println("Maximum number of vertex attributes supported: ", nrAttributes)

	r.createShaders()

	// Uncomment this call to draw in wireframe polygons.
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (r *Renderer) Teardown() {
	for _, m := range r.meshes {
		m.TearDown()
	}
}

func (r *Renderer) CreateMesh() *Mesh {
	m := &Mesh{}

	m.Setup()

	r.meshes = append(r.meshes, m)

	return m
}

func (r *Renderer) Draw(c *entity.Camera) {
	gl.UseProgram(r.shaderProgram)

	model := mgl32.Translate3D(0, 0, 0)
	location := gl.GetUniformLocation(r.shaderProgram, gl.Str("model\x00"))
	if location == -1 {
		panic("Could not get model location")
	}
	gl.UniformMatrix4fv(location, 1, false, &model[0])

	projection := mgl32.Perspective(45.0, float32(1440)/900, 0.1, 1000)
	location = gl.GetUniformLocation(r.shaderProgram, gl.Str("projection\x00"))
	if location == -1 {
		panic("Could not get projection location")
	}
	gl.UniformMatrix4fv(location, 1, false, &projection[0])

	view := c.CreateViewMatrix()
	location = gl.GetUniformLocation(r.shaderProgram, gl.Str("view\x00"))
	if location == -1 {
		panic("Could not get view location")
	}
	gl.UniformMatrix4fv(location, 1, false, &view[0])

	for i := 0; i < len(r.meshes); i++ {
		if r.meshes[i].active == false {
			r.meshes = append(r.meshes[:i], r.meshes[i+1:]...)
		}
	}

	for _, m := range r.meshes {
		if m.vertexCount > 0 {
			gl.BindVertexArray(m.vao)
			gl.DrawElements(gl.TRIANGLES, int32(len(m.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
	}
}

func (r *Renderer) createShaders() {
	{
		r.vertexShader = gl.CreateShader(gl.VERTEX_SHADER)

		shaderSrc, free := gl.Strs(vertexShader)
		defer free()

		gl.ShaderSource(r.vertexShader, 1, shaderSrc, nil)
		gl.CompileShader(r.vertexShader)

		var success int32
		gl.GetShaderiv(r.vertexShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(r.vertexShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(r.vertexShader, logLen, nil, &infoLog[0])

			panic(string(infoLog))
		}

		defer gl.DeleteShader(r.vertexShader)
	}

	{
		r.fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)

		shaderSrc, free := gl.Strs(fragmentShader)
		defer free()

		gl.ShaderSource(r.fragmentShader, 1, shaderSrc, nil)
		gl.CompileShader(r.fragmentShader)

		var success int32
		gl.GetShaderiv(r.fragmentShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(r.fragmentShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(r.fragmentShader, logLen, nil, &infoLog[0])

			panic(string(infoLog))
		}

		defer gl.DeleteShader(r.fragmentShader)
	}

	{
		r.shaderProgram = gl.CreateProgram()
		gl.AttachShader(r.shaderProgram, r.vertexShader)
		gl.AttachShader(r.shaderProgram, r.fragmentShader)
		gl.LinkProgram(r.shaderProgram)

		var success int32
		gl.GetProgramiv(r.shaderProgram, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(r.shaderProgram, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(r.shaderProgram, logLen, nil, &infoLog[0])

			panic(string(infoLog))
		}
	}
}
