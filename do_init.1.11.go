// +build go1.11
// +build !go1.13

package eleplugin

import "unsafe"

// Suitable for go 1.11 to 1.12

func (p *Plugin) callInit() error {
	initFuncPC, err := p.lookup(p.pkgPath + ".init")
	if err != nil {
		return err
	}
	initFuncP := &initFuncPC
	initFunc := *(*func())(unsafe.Pointer(&initFuncP))
	initFunc()
	return nil
}
