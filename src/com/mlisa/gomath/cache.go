package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type CacheManager struct {
	db *bolt.DB
}

func (cache *CacheManager) addNewOperation(operation string, result string) {
	var err error
	cache.db, err = bolt.Open("cacheDB", 0600, nil)
	if err != nil {
		log.Fatal("[ERROR]")
	}
	defer cache.db.Close()

	cache.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("OperationBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		b.Put([]byte(operation), []byte(result))
		return nil
	})
}

func (cache *CacheManager) retrieveResult(operation string) (string, bool) {
	var err error
	var result []byte
	cache.db, err = bolt.Open("cacheDB", 0600, nil)
	if err != nil {
		log.Fatal("[ERROR]")
	}
	defer cache.db.Close()

	err = cache.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("OperationBucket"))
		result = b.Get([]byte(operation))
		if result == nil {
			return bolt.ErrInvalid
		}
		return nil
	})

	if err != nil {
		return "", false
	}

	return string(result), true
}
