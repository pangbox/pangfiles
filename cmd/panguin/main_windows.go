// +build windows

package main

import "syscall"

func init() {
	dll := syscall.MustLoadDLL("user32")
	if proc, err := dll.FindProc("SetProcessDpiAwarenessContext"); err == nil {
		aware := -4
		proc.Call(uintptr(aware))
	}
}
