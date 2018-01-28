package processor

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	contextKey   = 1
	contextValue = contextKey
)

// errClosed error returned when we are closed
// i.e. we've been notified not to deal with any more input.
var errClosed = errors.New("closed")

// regexFilter regex to match non-alphanumeric characters.
var regexFilter = regexp.MustCompile(`[^\w]`)

type (
	// Context an alias to context.Context.
	Context = context.Context
	// cancelFunc an alias to context.CancelFunc
	cancelFunc = context.CancelFunc
)

// spawnContext spawns a context from the provided parent `ctx`
func spawnContext(ctx Context) Context {
	return context.WithValue(ctx, contextKey, contextValue)
}

// New initialises a new Processor with a context.
func New(ctx Context) Processor {
	p := &processor{store: &store{}}
	p.ctx, p.cancel = context.WithCancel(ctx)

	return p
}

// stat a wrapper for a data channel to push to
// and an instance of a Statistic.
type stat struct {
	data      chan<- Word
	statistic Statistic
}

// store an internal store of stats
// and open data channels.
type store struct {
	sync.Mutex
	stats []stat
}

// set sets a new stat to the store.
func (s *store) set(ch chan<- Word, st Statistic) {
	s.Lock()
	defer s.Unlock()

	s.stats = append(s.stats, stat{data: ch, statistic: st})
}

// push pushes a word to listening stats.
func (s *store) push(w Word) {
	s.Lock()
	defer s.Unlock()
	for _, stat := range s.stats {
		stat.data <- w
	}
}

// close closes all of the listening stats and attached
// data channels.
func (s *store) close() {
	s.Lock()
	defer s.Unlock()
	for _, stat := range s.stats {
		stat.statistic.Close()
		close(stat.data)
	}
}

// read writes data to the provided json encoder.
func (s *store) read(e *json.Encoder) error {
	s.Lock()
	defer s.Unlock()

	data := make(map[string]interface{})
	for _, stat := range s.stats {
		st := stat.statistic
		data[st.Name()] = st.Retrieve()
	}

	return e.Encode(data)
}

// Processor interface to a processor.
// Using an interface limits the functionality exposed outside
// of the package.
type Processor interface {
	io.Closer
	// Register registers a new statistic with the processor.
	Register(Statistic) error
	// Process processes an input stream with the registered
	// statistics.
	// Because we are utilising go's io.ReadCloser interface
	// input could come from any source which is a Reader. i.e.
	// a file, http input etc.
	Process(io.Reader) error
	// Write writes data to the provided writer from the
	// registered statistics.
	Write(io.Writer) error
}

// processor internal processor instance.
type processor struct {
	sync.RWMutex
	closed bool
	store  *store
	cancel cancelFunc
	ctx    Context
}

// isClosed checks to see if the processor is closed.
func (p *processor) isClosed() bool {
	p.Lock()
	defer p.Unlock()

	return p.closed
}

// Close implements io.Closer
func (p *processor) Close() error {
	if p.isClosed() {
		return errClosed
	}

	// cancel the attached context,
	// this will cancel any attached
	// statistics listening for input.
	p.cancel()
	p.store.close()
	p.closed = true

	return nil
}

// Register implements Processor interface.
func (p *processor) Register(s Statistic) error {
	input := s.Listen(spawnContext(p.ctx))
	p.store.set(input, s)

	return nil
}

// Process implements Processor interface.
func (p *processor) Process(r io.Reader) error {
	if p.isClosed() {
		return errClosed
	}

	p.Lock()
	defer p.Unlock()

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	var wg sync.WaitGroup

	// Scan all words from the input.
	for scanner.Scan() {
		wg.Add(1)

		word := regexFilter.ReplaceAllString(
			strings.ToLower(scanner.Text()), "",
		)

		go func(w Word) {
			defer wg.Done()
			p.store.push(w)
		}(Word(word))
	}

	wg.Wait()

	return nil
}

// Write implements Processor interface.
func (p *processor) Write(w io.Writer) error {
	p.Lock()
	defer p.Unlock()

	return p.store.read(json.NewEncoder(w))
}
