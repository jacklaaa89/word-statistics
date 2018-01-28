package statistics

const (
	nameTopFive = "top_5_words"
)

// NewTopFive initialises a new top five metric.
func NewTopFive() *TopFive {
	return &TopFive{baseCount: newBaseCount()}
}

// TopFive stat to return the top 5 seen words.
type TopFive struct {
	baseCount
}

// Name implements Statistic interface.
func (TopFive) Name() string {
	return nameTopFive
}
