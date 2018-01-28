package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// testStatName the name of the test stat.
	testStatName = "test_stat"
	// sampleInput sample input to use in tests.
	sampleInput = "This Is Four Words"
)

// testStatistic a statistic to use in tests.
type testStatistic struct {
	sync.Mutex
	closed bool
	words  []Word
}

// Close implements io.Closer
func (t *testStatistic) Close() error {
	t.closed = true
	return nil
}

// Listen implements Statistic interface.
func (t *testStatistic) Listen(ctx Context) chan<- Word {
	ch := make(chan Word, 10)
	go func() {
		for w := range ch {
			select {
			case <-ctx.Done():
				return
			default:
				t.Lock()
				t.words = append(t.words, w)
				t.Unlock()
			}
		}
	}()

	return ch
}

// Retrieve implements Statistic interface.
func (t *testStatistic) Retrieve() interface{} {
	t.Lock()
	defer t.Unlock()

	return t.words
}

// Name implements Statistic interface.
func (testStatistic) Name() string {
	return testStatName
}

// TestNew tests we can initialise a new processor.
func TestNew(t *testing.T) {
	p := New(context.Background())
	assert.IsType(t, &processor{}, p)
}

// TestProcessor_Process tests that we can process input.
func TestProcessor_Process(t *testing.T) {
	response := make(map[string][]string)

	p := New(context.Background())
	stat := &testStatistic{}

	p.Register(stat)
	p.Process(strings.NewReader(sampleInput))

	// we need to sleep at this point as
	// the process is done in the background
	// which is preferred on a web-server
	// (we can queue words to be processed).
	// but in this test it finishes prior to
	// processing all of the streamed input.
	// so a small sleep is required before the
	// read to ensure all data is consumed.
	time.Sleep(10 * time.Millisecond)
	b := bytes.NewBufferString("")
	p.Write(b)

	err := json.Unmarshal(b.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response, 1)
	words, ok := response[testStatName]
	assert.True(t, ok)
	assert.Len(t, words, 4)
}

// TestProcessor_Process_OnClosedProcessor tests that the
// processor will not accept any more input when it is closed.
func TestProcessor_Process_OnClosedProcessor(t *testing.T) {
	p := New(context.Background())
	stat := &testStatistic{}

	p.Register(stat)
	p.Close()
	assert.Error(t, p.Process(strings.NewReader(sampleInput)))
}

// TestProcessor_Close_OnClosedProcessor tests you cannot close a closed
// processor.
func TestProcessor_Close_OnClosedProcessor(t *testing.T) {
	p := New(context.Background())
	stat := &testStatistic{}

	p.Register(stat)
	p.Close()
	assert.Error(t, p.Close())
}

func TestProcessor_Process_WithNonAlphanumericChars(t *testing.T) {
	const text = `Text!!.£££###`

	response := make(map[string][]string)

	p := New(context.Background())
	stat := &testStatistic{}

	p.Register(stat)
	p.Process(strings.NewReader(text))

	time.Sleep(10 * time.Millisecond)
	b := bytes.NewBufferString("")
	p.Write(b)

	err := json.Unmarshal(b.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response, 1)
	words, ok := response[testStatName]
	assert.True(t, ok)
	assert.Len(t, words, 1)
	assert.Equal(t, `text`, words[0])
}
