package statistics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCount_Name(t *testing.T) {
	count := &Count{}
	assert.Equal(t, nameCount, count.Name())
}

func TestCount_Listen(t *testing.T) {
	count := &Count{}
	defer count.Close()

	data := count.Listen(context.Background())
	for i := 0; i < 5; i++ {
		text := sampleWord2
		if i%2 != 0 {
			text = sampleWord1
		}
		data <- word(text)
	}

	// ensure that we have processed the words.
	time.Sleep(10 * time.Millisecond)
	result, ok := count.Retrieve().(int)
	assert.True(t, ok)
	assert.Equal(t, result, 5)
}
