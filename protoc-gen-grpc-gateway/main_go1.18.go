//go:build go1.18

package main

import "runtime/debug"

func init() {
	v, ok := debug.ReadBuildInfo()
	if ok {
		version = v.Main.Version
	}
}
