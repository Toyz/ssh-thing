//go:build !windows
// +build !windows

package util

func EnableVirtualTerminalProcessing() error {
	return nil
}

func IsWindows() bool {
	return false
}
