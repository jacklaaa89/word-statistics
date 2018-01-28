package statistics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sampleWord1 = "test"
	sampleWord2 = "word"
)

func TestNewTopFive(t *testing.T) {
	topFive := NewTopFive()
	assert.IsType(t, &TopFive{}, topFive)
}

func TestTopFive_Name(t *testing.T) {
	assert.Equal(t, nameTopFive, NewTopFive().Name())
}

func TestTopFive_Listen(t *testing.T) {
	topFive := NewTopFive()
	data := topFive.Listen(context.Background())
	for i := 0; i < 5; i++ {
		text := sampleWord2
		if i%2 != 0 {
			text = sampleWord1
		}
		data <- word(text)
	}

	// ensure that we have processed the words.
	time.Sleep(10 * time.Millisecond)
	result, ok := topFive.Retrieve().(pairs)
	assert.True(t, ok)
	assert.Len(t, result, 2)
}
