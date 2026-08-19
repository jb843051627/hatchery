package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

func TestBug03_CacheUpdateConcurrentSafe(t *testing.T) {
	rc := NewReadingCache()
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			rc.Update(&model.SensorReading{
				IncubatorID: int64(n % 20),
				SensorType:  model.SensorTypeTemperature,
				Value:       float64(n),
				RecordedAt:  time.Now(),
			})
		}(i)
	}
	wg.Wait()
	snap := rc.Snapshot()
	if len(snap) == 0 {
		t.Fatal("expected non-empty cache after concurrent updates")
	}
}
