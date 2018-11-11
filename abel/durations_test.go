package abel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentAggregateWindow(t *testing.T) {
	now := time.Date(2018, 11, 11, 16, 00, 00, 00, time.UTC)
	expected := time.Date(2018, 11, 11, 00, 00, 00, 00, time.UTC)
	assert.Equal(t, expected.Unix()*1000, CurrentAggregateWindow(86400000, now))
}

func TestCurrentAggregateWindowWithZeroGranularity(t *testing.T) {
	now := time.Date(2018, 11, 11, 16, 00, 00, 00, time.UTC)
	assert.Equal(t, int64(-1), CurrentAggregateWindow(0, now))
}

func TestPreviousAggregateWindow(t *testing.T) {
	now := time.Date(2018, 11, 11, 16, 00, 00, 00, time.UTC)
	expected := time.Date(2018, 11, 10, 00, 00, 00, 00, time.UTC)
	assert.Equal(t, expected.Unix()*1000, PreviousAggregateWindow(86400000, now))
}

func TestNextAggregateWindow(t *testing.T) {
	now := time.Date(2018, 11, 11, 16, 00, 00, 00, time.UTC)
	expected := time.Date(2018, 11, 12, 00, 00, 00, 00, time.UTC)
	assert.Equal(t, expected.Unix()*1000, NextAggregateWindow(86400000, now))
}

func TestTimeToNextAggregateWindow(t *testing.T) {
	now := time.Date(2018, 11, 11, 16, 00, 00, 00, time.UTC)
	assert.Equal(t, time.Duration(8*time.Hour), TimeToNextAggregateWindow(86400000, now))
}
