package cache

import (
	"sync"

	"github.com/sangnt1552314/digimontex/internal/models"
)

type DigimonCache struct {
	data  map[int]models.DigimonDetail
	order []int
	mutex sync.RWMutex
	size  int
}

func NewDigimonCache(size int) *DigimonCache {
	return &DigimonCache{
		data:  make(map[int]models.DigimonDetail),
		order: make([]int, 0, size),
		size:  size,
	}
}

func (c *DigimonCache) Get(id int) (*models.DigimonDetail, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	digimon, exists := c.data[id]
	if exists {
		// Move to front when accessed (LRU behavior)
		c.moveToFrontUnsafe(id)
	}
	return &digimon, exists
}

func (c *DigimonCache) Put(id int, digimon *models.DigimonDetail) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If already exists, update and move to front
	if _, exists := c.data[id]; exists {
		c.moveToFrontUnsafe(id)
		c.data[id] = *digimon
		return
	}

	// If cache is full, remove oldest (first in order)
	if len(c.order) >= c.size {
		oldest := c.order[0]
		delete(c.data, oldest)
		c.order = c.order[1:]
	}

	// Add new entry to the end (most recent)
	c.data[id] = *digimon
	c.order = append(c.order, id)
}

func (c *DigimonCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[int]models.DigimonDetail)
	c.order = c.order[:0]
}

func (c *DigimonCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.data)
}

func (c *DigimonCache) GetRecentIDs() []int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Return a copy of the order slice (most recent last)
	result := make([]int, len(c.order))
	copy(result, c.order)
	return result
}

// moveToFrontUnsafe moves an existing ID to the end of order slice (most recent)
// This method assumes the mutex is already locked
func (c *DigimonCache) moveToFrontUnsafe(id int) {
	for i, v := range c.order {
		if v == id {
			// Remove from current position
			c.order = append(c.order[:i], c.order[i+1:]...)
			// Add to end (most recent)
			c.order = append(c.order, id)
			break
		}
	}
}
