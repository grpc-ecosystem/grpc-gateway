//go:build !go1.18

package main

func readVersion() string {
	return version
}
