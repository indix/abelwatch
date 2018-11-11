package abel

import "time"

// CurrentAggregateWindow returns the current timestamp (in millis) at which the aggregate is happening
func CurrentAggregateWindow(granularity int64, now time.Time) int64 {
	if granularity == int64(0) {
		return int64(-1)
	}
	return (now.Unix() * 1000) / granularity * granularity
}

// PreviousAggregateWindow returns the previous timestamp (in millis) at which the aggregate happened
func PreviousAggregateWindow(granularity int64, now time.Time) int64 {
	return CurrentAggregateWindow(granularity, now) - granularity
}

// NextAggregateWindow returns the next timestamp at which the aggregate will happen
func NextAggregateWindow(granularity int64, now time.Time) int64 {
	return CurrentAggregateWindow(granularity, now) + granularity
}

// TimeToNextAggregateWindow returns the duration to next aggregate given now as the current time
func TimeToNextAggregateWindow(granularity int64, now time.Time) time.Duration {
	currentWindow := CurrentAggregateWindow(granularity, now)
	return time.Millisecond * time.Duration(granularity-((now.Unix()*1000)-currentWindow))
}
