package statistics

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"sync"

	"word-statistics/stats/processor"
)

const (
	wordBuffer = 10
)

// Context an alias to the internal context package.
type Context = context.Context

// pairs a list of pairs.
type pairs []pair

// Less implements sort.Interface interface.
func (p pairs) Less(i, j int) bool {
	return p[i].frequency > p[j].frequency
}

// Swap implements sort.Interface interface.
func (p pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Len implements sort.Interface interface.
func (p pairs) Len() int {
	return len(p)
}

// pair represents a word / frequency seen pair when sorting.
type pair struct {
	word      word
	frequency int
}

// MarshalJSON implements json.Marshaller interface.
func (p pair) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p.word))
}

// word an alias to processor.Word
type word = processor.Word

// action an action to apply to each seen word.
type action = func(w word)

// errClosed error returned when we are closed
// i.e. we've been notified not to deal with any more input.
var errClosed = errors.New("closed")

// statistic base implementation of a statistic.
type statistic struct {
	sync.Mutex
	closed bool
}

// isClosed checks to see if this stat is closed.
func (s *statistic) isClosed() bool {
	s.Lock()
	defer s.Unlock()

	return s.closed
}

// Name implements io.Closer interface.
func (s *statistic) Close() error {
	if s.isClosed() {
		return errClosed
	}

	s.closed = true
	return nil
}

// listen helper func which handles the state of the returned channel.
// applies the supplied action for every word received.
func (s *statistic) listen(ctx Context, action action) chan<- word {
	ch := make(chan word, wordBuffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Close()
				close(ch)
				return
			case word := <-ch:
				// increment the counter if the
				// stat is not closed.
				if s.isClosed() {
					return
				}
				action(word)
			}
		}
	}()

	return ch
}

// top gets the top pairs in pairs up to `amount` in length
func (statistic) top(pairs pairs, amount int) pairs {
	// sort by frequency desc to get the most frequent words in the pairs list.
	sort.Sort(pairs)

	// reduce the pairs to the required amount.
	if amount > len(pairs) {
		amount = len(pairs)
	}

	pairs = pairs[:amount]

	// apply a alphanumeric sort so ensure the order of the returned
	// list.
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].word < pairs[j].word
	})

	return pairs
}
