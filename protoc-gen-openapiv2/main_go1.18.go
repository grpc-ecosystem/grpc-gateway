//go:build go1.18

package main

import "runtime/debug"

func readVersion() string {
	v, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	return v.Main.Version
}
