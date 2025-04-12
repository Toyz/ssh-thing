//go:build windows
// +build windows

package util

import (
	"os"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	ansiEnabled     bool
	ansiEnableMutex sync.Mutex
)

func EnableVirtualTerminalProcessing() error {
	ansiEnableMutex.Lock()
	defer ansiEnableMutex.Unlock()

	if ansiEnabled {
		return nil
	}

	stdout := windows.Handle(os.Stdout.Fd())

	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err != nil {
		return err
	}

	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if err := windows.SetConsoleMode(stdout, mode); err != nil {
		return err
	}

	stderr := windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(stderr, &mode); err != nil {
		return nil
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	windows.SetConsoleMode(stderr, mode)

	ansiEnabled = true
	return nil
}

func IsWindows() bool {
	return os.Getenv("OS") == "Windows_NT" || os.Getenv("GOOS") == "windows"
}
