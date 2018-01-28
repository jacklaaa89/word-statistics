package statistics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTopLetters(t *testing.T) {
	topLetters := NewTopLetters()
	assert.IsType(t, &TopLetters{}, topLetters)
}

func TestTopLetters_Name(t *testing.T) {
	assert.Equal(t, nameTopLetters, NewTopLetters().Name())
}

func TestTopLetters_Listen(t *testing.T) {
	topLetters := NewTopLetters()
	data := topLetters.Listen(context.Background())
	for i := 0; i < 5; i++ {
		text := sampleWord2
		if i%2 != 0 {
			text = sampleWord1
		}
		data <- word(text)
	}

	// ensure that we have processed the words.
	time.Sleep(10 * time.Millisecond)
	result, ok := topLetters.Retrieve().(pairs)
	assert.True(t, ok)
	assert.Len(t, result, 5)
}
