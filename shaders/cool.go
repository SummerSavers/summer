package shaders

import "github.com/Noofbiz/pixelshader"

var CoolShader = &pixelshader.PixelShader{FragShader: `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

uniform vec2 u_resolution;  // Canvas size (width,height)
uniform vec2 u_mouse;       // mouse position in screen pixels
uniform float u_time;       // Time in seconds since load

float circularOut(float t) {
  return sqrt((1.0 - t) * t);
}

float circle(in vec2 _st, in float _radius){
    vec2 dist = _st-vec2(0.5);
	return 1.-smoothstep(_radius-(_radius*0.2),
                         _radius+(_radius*0.01),
                         dot(dist,dist)*4.0);
}

vec3 colorA = vec3(0.529, 0.808, 0.922);
vec3 colorB = vec3(0.4273,0.3806,0.1920);
vec3 white = vec3(0.90, 0.90, 0.90);

void main() {
    vec2 st = gl_FragCoord.xy/u_resolution.xy;

    float t = u_time*0.5;
    float x = u_mouse.x / u_resolution.x;
    st.x += 0.5 - x;
    st.y *= abs(sin(t));
    float pct = circularOut(st.y);

    gl_FragColor = vec4(vec3(mix(vec3(mix(colorA, colorB, pct)), white, circle(st,0.3))),1.0);
}
  `}
