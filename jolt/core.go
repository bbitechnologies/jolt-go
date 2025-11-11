package jolt

// #include "wrapper/core.h"
import "C"
import "fmt"

// Init initializes Jolt Physics (call once at startup)
func Init() error {
	result := C.JoltInit()
	if result == 0 {
		return fmt.Errorf("failed to initialize Jolt")
	}
	return nil
}

// Shutdown cleans up Jolt resources (call once at exit)
func Shutdown() {
	C.JoltShutdown()
}
