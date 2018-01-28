package statistics

import "word-statistics/stats/processor"

const (
	nameCount = "count"
)

// Count a statistic which counts all of the words which
// come in.
type Count struct {
	statistic
	counter int
}

// increment increments the counter by one.
func (c *Count) increment() {
	c.Lock()
	defer c.Unlock()

	c.counter++
}

// Listen implements Statistic interface.
func (c *Count) Listen(ctx processor.Context) chan<- word {
	return c.listen(ctx, func(w word) {
		c.increment()
	})
}

// Retrieve implements Statistic interface.
func (c *Count) Retrieve() interface{} {
	c.Lock()
	defer c.Unlock()

	return c.counter
}

// Name implements Statistic interface.
func (Count) Name() string {
	return nameCount
}
