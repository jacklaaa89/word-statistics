package processor

import "io"

// Word a representation of a `word` received as input.
type Word string

// Statistic a representation of a statistic.
// A statistic maintains its own state. Reacts to input
// provided on Listen and produces the current statistic
// with Retrieve.
// Statistics should be implemented thread-safe as we may retrieve and push
// input concurrently.
type Statistic interface {
	io.Closer
	// Listen returns a channel to listen on
	// for input. A buffered channel would increase performance.
	Listen(Context) chan<- Word
	// Retrieve retrieves this statistics value.
	Retrieve() interface{}
	// Name returns the name of the statistic.
	Name() string
}
