package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	expectedFrequency      = 1
	expectedOutputTemplate = `"%v"`
)

func TestPair_MarshalJSON(t *testing.T) {
	pair := pair{word: sampleWord1, frequency: expectedFrequency}
	d, err := json.Marshal(pair)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf(expectedOutputTemplate, sampleWord1), string(d))
}

func TestStatistic_Close_OnClosedStatistic(t *testing.T) {
	stat := &Count{}
	assert.NoError(t, stat.Close())
	assert.Error(t, stat.Close())
}

func TestStatistic_ListenOnCancelledContext(t *testing.T) {
	stat := &Count{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stat.Listen(ctx)

	time.Sleep(10 * time.Millisecond)
	assert.Error(t, stat.Close())
}

func TestStatistic_Listen_OnClosedStatistic(t *testing.T) {
	stat := &Count{}
	stat.Close()

	data := stat.Listen(context.Background())
	data <- word(sampleWord1)

	time.Sleep(10 * time.Millisecond)
	assert.Error(t, stat.Close())
}
