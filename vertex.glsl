#version 330

in vec3 vertex_position;
uniform mat4 proj, view, model;

void main() {
        gl_Position = proj * view * model * vec4(vertex_position, 1);
}
