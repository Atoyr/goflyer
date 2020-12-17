package db

import (
	"time"

	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/boltdb/bolt"
)

type Bolt struct {
	dbFile string
}

const (
	executionBucket = "Execution"
)

func NewBolt(dbFile string) (*Bolt, error) {
	b := new(Bolt)
	b.dbFile = dbFile
	return b, nil
}

func (b *Bolt) init() error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(tickerBucketName))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(executionBucketName))
		if err != nil {
			return err
		}
		for _, v := range models.Durations() {
			dbucket, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return err
			}
			_, err = dbucket.CreateBucketIfNotExists([]byte("open"))
			if err != nil {
				return err
			}
			_, err = dbucket.CreateBucketIfNotExists([]byte("close"))
			if err != nil {
				return err
			}
			_, err = dbucket.CreateBucketIfNotExists([]byte("high"))
			if err != nil {
				return err
			}
			_, err = dbucket.CreateBucketIfNotExists([]byte("low"))
			if err != nil {
				return err
			}
			_, err = dbucket.CreateBucketIfNotExists([]byte("volume"))
			if err != nil {
				return err
			}
		}
		_, err = tx.CreateBucketIfNotExists([]byte(logBucketName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Bolt) db() *bolt.DB {
	db, err := bolt.Open(b.dbFile, 0600, nil)
	if err != nil {
		logf("%v", err)
		return nil
	}
	return db
}

func (b *Bolt) InsertExecution(executions []bitflyer.Execution) {

}

func (b *Bolt) SelectExecution(after time.Time, count int) []bitflyer.Execution {
	return nil
}
