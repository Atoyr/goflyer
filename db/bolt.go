package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/atoyr/goflyer/models"
	"github.com/boltdb/bolt"
)

type Bolt struct {
	dbFile string
}

// Bucket layout
// - TickerBucket
// - DurationBucket
//   - CandleBucket
//   - SmasBucket
//   - EmasBucket

const (
	TickerBucket   = "Ticker"
	DurationBucket = "Duration"
	CandleBucketk  = "Candle"
	SmasBucket     = "Smas"
	EmasBucket     = "Emas"
)

func GetBolt(dbFile string) (Bolt, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	var b Bolt
	if err != nil {
		return b, err
	}
	defer db.Close()
	b.dbFile = dbFile
	b.init()
	return b, nil
}

func (b *Bolt) db() *bolt.DB {
	db, err := bolt.Open(b.dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func (b *Bolt) init() error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TickerBucket))
		if err != nil {
			return err
		}
		durationBucket, err := tx.CreateBucketIfNotExists([]byte(DurationBucket))
		if err != nil {
			return err
		}
		_, err = durationBucket.CreateBucketIfNotExists([]byte(CandleBucketk))
		if err != nil {
			return err
		}
		_, err = durationBucket.CreateBucketIfNotExists([]byte(SmasBucket))
		if err != nil {
			return err
		}
		_, err = durationBucket.CreateBucketIfNotExists([]byte(EmasBucket))
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

func (b *Bolt) UpdateTicker(t models.Ticker) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TickerBucket))
		marshalID := []byte(t.TickID)
		if buf, err := json.Marshal(t); err != nil {
			return err
		} else if err = bucket.Put(marshalID, buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Bolt) GetTicker(tickerID float64) (models.Ticker, error) {
	db := b.db()
	defer db.Close()
	ticker := new(models.Ticker)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(TickerBucket)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			min, err := tickerID
			if err != nil {
				return err
			}
			max, err := tickerID
			if err != nil {
				return err
			}

			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				ticker, err = models.JsonUnmarshalTicker(v)
				if err != nil {
					return err
				}
				fmt.Printf("%s: %s\n", k, v)
			}
			return nil
		}
		return nil
	})

	return *ticker, err
}

func (b *Bolt) GetTickerAll() ([]models.Ticker, error) {
	db := b.db()
	defer db.Close()
	tickers := make([]models.Ticker, 0)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(TickerBucket)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t models.Ticker
			json.Unmarshal(v, &t)
			tickers = append(tickers, t)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tickers, nil
}

func (b *Bolt) UpdateCandle(c models.Candle) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucketName := fmt.Sprintf("Candle_%s", c.Key())
		bucket := tx.Bucket([]byte(bucketName))
		if buf, err := json.Marshal(c); err != nil {
			return err
		} else if err = bucket.Put([]byte(c.Key()), buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (b *Bolt) GetCandleCollection() {

}
