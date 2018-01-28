package statistics

const nameTopLetters = "top_5_letters"

// NewTopLetters initialises a new top five metric.
func NewTopLetters() *TopLetters {
	return &TopLetters{baseCount: newBaseCount()}
}

// TopLetters stat to return the top 5 used letters.
type TopLetters struct {
	baseCount
}

// Listen implements Statistic interface.
func (t *TopLetters) Listen(ctx Context) chan<- word {
	return t.listen(ctx, t.process)
}

// process increments a counter of all of the letters
// seen. Any character that is not a letter is ignored.
func (t *TopLetters) process(w word) {
	for _, char := range w {
		t.baseCount.process(word(char))
	}
}

// Name implements Statistic interface.
func (TopLetters) Name() string {
	return nameTopLetters
}
