// +build linux,cgo

package eleplugin

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

#include <stdio.h>

static void* pluginOpen(const char* path, char** err) {
	// void* h = dlopen(path, RTLD_LAZY);
	void* h = dlopen(path, RTLD_NOW| RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return h;
}

static void* pluginLookup(void *h, const char* name, char** err) {
	void* r = dlsym(h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}

*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var pluginsMu sync.Mutex
var plugins map[string]*Plugin

type Plugin struct {
	filename string
	pkgPath string
	h        unsafe.Pointer
}

func doFuncPC(pc unsafe.Pointer) {
	p := &pc
	f := *(*func())(unsafe.Pointer(&p))
	f()
}


func Open(filename, pkgPath string) (*Plugin, error) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	if p, ok := plugins[filename]; ok {
		return p, nil
	}

	cFilename := C.CString(filename)
	var cErr *C.char
	h := C.pluginOpen(cFilename, &cErr)
	C.free(unsafe.Pointer(cFilename))

	if h == nil {
		return nil, fmt.Errorf("open plugin %q failed: %v", filename, C.GoString(cErr))
	}

	lastmoduleinit()

	p := &Plugin{
		filename: filename,
		pkgPath: pkgPath,
		h:        h,
	}

	err := p.callInit()
	if err != nil {
		return nil, err
	}

	if plugins == nil {
		plugins = make(map[string]*Plugin)
	}
	plugins[filename] = p
	return p, nil
}


func (p *Plugin) lookup(name string) (unsafe.Pointer, error) {
	cName := C.CString(name)
	var cErr *C.char
	ptr := C.pluginLookup(p.h, cName, &cErr)
	C.free(unsafe.Pointer(cName))

	if ptr == nil {
		return nil, fmt.Errorf("can not find symbol %q: %v", name, C.GoString(cErr))
	}
	return ptr, nil
}

func (p *Plugin) Start(name string) error {
	if name == "" {
		name = "PluginStart"
		// TODO go >= 1.17 need to add .abiinternal
	}
	f, err := p.lookup(p.pkgPath + "." + name)
	if err != nil {
		return err
	}
	doFuncPC(f)
	return nil
}
