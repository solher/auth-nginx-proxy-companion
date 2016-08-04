package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/solher/auth-nginx-proxy-companion/models"

	"github.com/boltdb/bolt"
)

type (
	GarbageCollectorSessionsRepo interface {
		Update(func(tx *bolt.Tx) error) error
	}

	GarbageCollector struct {
		repo GarbageCollectorSessionsRepo
	}
)

func NewGarbageCollector(repo GarbageCollectorSessionsRepo) *GarbageCollector {
	return &GarbageCollector{repo: repo}
}

func (gc *GarbageCollector) Run(dbLocation string, freq time.Duration) error {
	db, err := bolt.Open(dbLocation, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}

	go gc.run(db, freq)

	return nil
}

func (gc *GarbageCollector) run(db *bolt.DB, freq time.Duration) {
	tick := time.Tick(freq)

	for range tick {
		fmt.Println("Running the garbage collection...")

		if err := gc.collect(db); err != nil {
			fmt.Println("WARNING: errors occured during the garbage collection")
			continue
		}

		fmt.Println("Done.")
	}
}

func (gc *GarbageCollector) collect(db *bolt.DB) error {
	return gc.repo.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("sessions")).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			session := models.Session{}

			if err := json.Unmarshal(v, &session); err != nil {
				return err
			}

			if session.ValidTo.After(time.Now()) {
				continue
			}

			err := db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("sessions"))
				if err != nil {
					return err
				}

				b.Put(k, v)

				return nil
			})

			if err != nil {
				return err
			}

			if err := c.Delete(); err != nil {
				return err
			}
		}

		return nil
	})
}
