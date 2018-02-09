package envsensor

import (
	"time"
)

type CachedReading struct {
	Reading
	expiresAt time.Time
}

// reading is the Reading you want to return.
// cacheDuration is how long to keep this reading fresh for
func NewCachedReading(reading Reading, cacheDuration time.Duration) CachedReading {
	return CachedReading{
		Reading:   reading,
		expiresAt: time.Now().Add(cacheDuration),
	}
}

// Whether we are a stale reading or not
func (c *CachedReading) IsStale() bool {
	return time.Now().After(c.expiresAt)
}
