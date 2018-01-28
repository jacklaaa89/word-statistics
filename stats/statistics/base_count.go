package statistics

const wordCount = 5

// newBaseCount initialises a new baseCount
func newBaseCount() baseCount {
	return baseCount{words: make(map[word]int)}
}

// baseCount a base structure for holding statistics
// on the top 5 words that it sees.
type baseCount struct {
	statistic
	words map[word]int
}

// process the default process method, increments
// the map entry for the received word.
func (b *baseCount) process(w word) {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.words[w]; !ok {
		b.words[w] = 0
	}

	b.words[w] += 1
}

// Listen implements Statistic interface.
func (b *baseCount) Listen(ctx Context) chan<- word {
	return b.listen(ctx, b.process)
}

// Retrieve implements Statistic interface.
func (b *baseCount) Retrieve() interface{} {
	b.Lock()
	defer b.Unlock()

	var index int
	pairs := make(pairs, len(b.words))
	for word, frequency := range b.words {
		pairs[index] = pair{word: word, frequency: frequency}
		index++
	}

	return b.top(pairs, wordCount)
}
