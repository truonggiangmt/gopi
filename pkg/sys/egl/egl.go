// +build egl

package egl

import (
	"unsafe"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: egl
#include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EGLDisplay      C.EGLDisplay
	EGLConfig       C.EGLConfig
	EGLContext      C.EGLContext
	EGLSurface      C.EGLSurface
	EGLNativeWindow unsafe.Pointer
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	EGLSurfaceTypeMap = map[string]gopi.SurfaceFlags{
		"OpenGL":     gopi.SURFACE_FLAG_OPENGL,
		"OpenGL_ES":  gopi.SURFACE_FLAG_OPENGL_ES,
		"OpenGL_ES2": gopi.SURFACE_FLAG_OPENGL_ES2,
		"OpenVG":     gopi.SURFACE_FLAG_OPENVG,
	}
	EGLAPIMap = map[gopi.SurfaceFlags]EGLAPI{
		gopi.SURFACE_FLAG_OPENGL_ES: EGL_API_OPENGL_ES,
		gopi.SURFACE_FLAG_OPENVG:    EGL_API_OPENVG,
		gopi.SURFACE_FLAG_OPENGL:    EGL_API_OPENGL,
	}
	EGLRenderableMap = map[gopi.SurfaceFlags]EGLRenderableFlag{
		gopi.SURFACE_FLAG_OPENGL:    EGL_RENDERABLE_FLAG_OPENGL,
		gopi.SURFACE_FLAG_OPENGL_ES: EGL_RENDERABLE_FLAG_OPENGL_ES,
		gopi.SURFACE_FLAG_OPENVG:    EGL_RENDERABLE_FLAG_OPENVG,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGLGetDisplay(display uint) EGLDisplay {
	return EGLDisplay(C.eglGetDisplay(C.EGLNativeDisplayType(uintptr(display))))
}

func EGLInitialize(display EGLDisplay) (int, int, error) {
	var major, minor C.EGLint
	if C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor))) != EGL_TRUE {
		return 0, 0, EGLGetError()
	} else {
		return int(major), int(minor), nil
	}
}

func EGLTerminate(display EGLDisplay) error {
	if C.eglTerminate(C.EGLDisplay(display)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLQueryString(display EGLDisplay, value EGLQuery) string {
	return C.GoString(C.eglQueryString(C.EGLDisplay(display), C.EGLint(value)))
}

////////////////////////////////////////////////////////////////////////////////
// SURFACE CONFIGS

func EGLGetConfigs(display EGLDisplay) ([]EGLConfig, error) {
	var num_config C.EGLint
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	if num_config == C.EGLint(0) {
		return nil, EGL_BAD_CONFIG
	}
	// configs is a slice so we need to pass the slice pointer
	configs := make([]EGLConfig, num_config)
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	} else {
		return configs, nil
	}
}

func EGLGetConfigAttrib(display EGLDisplay, config EGLConfig, attrib EGLConfigAttrib) (int, error) {
	var value C.EGLint
	if C.eglGetConfigAttrib(C.EGLDisplay(display), C.EGLConfig(config), C.EGLint(attrib), &value) != EGL_TRUE {
		return 0, EGLGetError()
	} else {
		return int(value), nil
	}
}

func EGLGetConfigAttribs(display EGLDisplay, config EGLConfig) (map[EGLConfigAttrib]int, error) {
	attribs := make(map[EGLConfigAttrib]int, 0)
	for k := EGL_COMFIG_ATTRIB_MIN; k <= EGL_COMFIG_ATTRIB_MAX; k++ {
		if v, err := EGLGetConfigAttrib(display, config, k); err == EGL_BAD_ATTRIBUTE {
			continue
		} else if err != nil {
			return nil, err
		} else {
			attribs[k] = v
		}
	}
	return attribs, nil
}

func EGLChooseConfig_(display EGLDisplay, attributes map[EGLConfigAttrib]int) ([]EGLConfig, error) {
	var num_config C.EGLint

	// Make list of attributes as eglInt values
	attribute_list := make([]C.EGLint, len(attributes)*2+1)
	i := 0
	for k, v := range attributes {
		attribute_list[i] = C.EGLint(k)
		attribute_list[i+1] = C.EGLint(v)
		i = i + 2
	}
	attribute_list[i] = C.EGLint(EGL_NONE)

	// Get number of configurations this matches
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	// Return EGL_BAD_ATTRIBUTE if the attribute set doesn't match
	if num_config == 0 {
		return nil, EGL_BAD_ATTRIBUTE
	}
	// Allocate an array
	configs := make([]EGLConfig, num_config)
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	// Return the configurations
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	} else {
		return configs, nil
	}
}

func EGLChooseConfig(display EGLDisplay, r_bits, g_bits, b_bits, a_bits uint, surface_type EGLSurfaceTypeFlag, renderable_type EGLRenderableFlag) (EGLConfig, error) {
	if configs, err := EGLChooseConfig_(display, map[EGLConfigAttrib]int{
		EGL_RED_SIZE:        int(r_bits),
		EGL_GREEN_SIZE:      int(g_bits),
		EGL_BLUE_SIZE:       int(b_bits),
		EGL_ALPHA_SIZE:      int(a_bits),
		EGL_SURFACE_TYPE:    int(surface_type),
		EGL_RENDERABLE_TYPE: int(renderable_type),
	}); err != nil {
		return 0, err
	} else if len(configs) == 0 {
		return 0, EGL_BAD_CONFIG
	} else {
		return configs[0], nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// API

func EGLBindAPI(api EGLAPI) error {
	if success := C.eglBindAPI(C.EGLenum(api)); success != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLQueryAPI() (EGLAPI, error) {
	if api := EGLAPI(C.eglQueryAPI()); api == 0 {
		return EGL_API_NONE, EGLGetError()
	} else {
		return api, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONTEXT

func EGLCreateContext(display EGLDisplay, config EGLConfig, share_context EGLContext) (EGLContext, error) {
	if context := EGLContext(C.eglCreateContext(C.EGLDisplay(display), C.EGLConfig(config), C.EGLContext(share_context), nil)); context == nil {
		return nil, EGLGetError()
	} else {
		return context, nil
	}
}

func EGL_DestroyContext(display EGLDisplay, context EGLContext) error {
	if C.eglDestroyContext(C.EGLDisplay(display), C.EGLContext(context)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLMakeCurrent(display EGLDisplay, draw, read EGLSurface, context EGLContext) error {
	if C.eglMakeCurrent(C.EGLDisplay(display), C.EGLSurface(draw), C.EGLSurface(read), C.EGLContext(context)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SURFACE

func EGLCreateSurface(display EGLDisplay, config EGLConfig, window EGLNativeWindow) (EGLSurface, error) {
	if surface := EGLSurface(C.eglCreateWindowSurface(C.EGLDisplay(display), C.EGLConfig(config), C.EGLNativeWindowType(window), nil)); surface == nil {
		return nil, EGLGetError()
	} else {
		return surface, nil
	}
}

func EGLDestroySurface(display EGLDisplay, surface EGLSurface) error {
	if C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLSwapBuffers(display EGLDisplay, surface EGLSurface) error {
	if C.eglSwapBuffers(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}