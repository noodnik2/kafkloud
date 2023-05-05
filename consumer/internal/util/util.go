package util

import "io"

// ComponentErrorHandler handles errors and returns true iff the component should die
type ComponentErrorHandler func(reason error) bool

type Component interface {
	
	// Launch the component; if this function returns nil, the "clean up function"
	// will be called if/when the background process that was launched completes
	Launch(ComponentErrorHandler) error

	// Closer for the component; signals the background process should stop at its
	// earliest convenience for a "graceful" shutdown
	io.Closer
}
