package cache

import (
	"sync"

	"github.com/jb843051627/hatchery/internal/model"
)

type ReadingCache struct {
	mu     sync.RWMutex
	latest map[int64]map[string]*model.SensorReading
}

func NewReadingCache() *ReadingCache {
	return &ReadingCache{latest: make(map[int64]map[string]*model.SensorReading)}
}

func (c *ReadingCache) Update(r *model.SensorReading) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.latest[r.IncubatorID] == nil {
		c.latest[r.IncubatorID] = make(map[string]*model.SensorReading)
	}
	c.latest[r.IncubatorID][r.SensorType] = r
}

func (c *ReadingCache) Get(incubatorID int64, sensorType string) (*model.SensorReading, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if m, ok := c.latest[incubatorID]; ok {
		if r, ok := m[sensorType]; ok {
			return r, true
		}
	}
	return nil, false
}

func (c *ReadingCache) Snapshot() map[int64]map[string]*model.SensorReading {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[int64]map[string]*model.SensorReading, len(c.latest))
	for k, v := range c.latest {
		nv := make(map[string]*model.SensorReading, len(v))
		for sk, sv := range v {
			cp := *sv
			nv[sk] = &cp
		}
		out[k] = nv
	}
	return out
}

func (c *ReadingCache) AllForIncubator(incubatorID int64) map[string]*model.SensorReading {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if m, ok := c.latest[incubatorID]; ok {
		out := make(map[string]*model.SensorReading, len(m))
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	return nil
}

func (c *ReadingCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latest = make(map[int64]map[string]*model.SensorReading)
}
