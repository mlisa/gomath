package peer

import (
	"errors"

	gocache "github.com/patrickmn/go-cache"
)

type CacheManager struct {
	c *gocache.Cache
}

func (cache *CacheManager) addNewOperation(operation string, result string) error {
	if cache.c == nil {
		cache.c = gocache.New(gocache.NoExpiration, gocache.DefaultExpiration)
	}
	if err := cache.c.Add(operation, result, gocache.NoExpiration); err != nil {
		return err
	}
	return nil
}

func (cache *CacheManager) retrieveResult(operation string) (string, error) {
	if cache.c != nil {
		if item, found := cache.c.Get(operation); found {
			return item.(string), nil
		} else {
			return "", nil
		}
	}
	return "", errors.New("[CACHE] Cache not found")
}
