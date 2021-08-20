// +build go1.13

package eleplugin

import "unsafe"

// // Suitable for go >= 1.13

func (p *Plugin) callInit() error {
	t, err := p.lookup(p.pkgPath + "..inittask")
	if err != nil {
		return err
	}
	doInit(t)
	println("after do Init")
	return nil
}

// doInit is defined in package runtime
//go:linkname doInit runtime.doInit
func doInit(t unsafe.Pointer) // t should be a *runtime.initTask