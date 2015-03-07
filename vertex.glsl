#version 330

uniform mat4 proj, view, model;
in vec3 vertex_position;
out vec3 frag_color;

void main() {
        gl_Position = proj * view * model * vec4(vertex_position, 1);
        frag_color = vertex_position;
}
