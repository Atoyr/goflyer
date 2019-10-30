package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
	"github.com/boltdb/bolt"
)

type Bolt struct {
	dbFile string
}

// Bucket layout
// - tickerBucket
// - DurationBucket
//   - CandleBucket
//   - SmasBucket
//   - EmasBucket

const (
	tickerBucket      = "Ticker"
	executionBucket   = "Execution"
	durationBucket    = "Duration"
	candleBucket      = "Candle"
	smasBucket        = "Smas"
	emasBucket        = "Emas"
)

func GetBolt(dbFile string) (Bolt, error) {
	var b Bolt
	b.dbFile = dbFile
	err := b.init()
	return b, err
}

func getCandleBucketName(duration int64) string{
	d := time.Duration(duration)
	return fmt.Sprintf("%s_%s",durationBucket,models.GetDurationString(d))
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
		_, err := tx.CreateBucketIfNotExists([]byte(tickerBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(executionBucket))
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
		bucket := tx.Bucket([]byte(tickerBucket))
		marshalID := util.Float64ToBytes(t.TickID)
		if buf, err := json.Marshal(t); err != nil {
			return err
		} else if err = bucket.Put(marshalID, buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (b *Bolt) GetTicker(tickID float64) (models.Ticker, error) {
	db := b.db()
	defer db.Close()
	ticker := new(models.Ticker)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(tickerBucket)).Cursor()

		marshalID := util.Float64ToBytes(tickID)

		for k, v := c.Seek(marshalID); k != nil && bytes.Compare(k, marshalID) <= 0; k, v = c.Next() {
			t, err := models.JsonUnmarshalTicker(v)
			if err != nil {
				return err
			}
			ticker = t
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
		c := tx.Bucket([]byte(tickerBucket)).Cursor()

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

func (b *Bolt) UpdateExecution(execution models.Execution) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(executionBucket))
		marshalID := util.Float64ToBytes(execution.ID)
		if buf, err := json.Marshal(execution); err != nil {
			return err
		} else if err = bucket.Put(marshalID, buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil 
}
func (b *Bolt) GetExecutionAll() ([]models.Execution, error) {
	db := b.db()
	defer db.Close()
	executions := make([]models.Execution, 0)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(executionBucket)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e models.Execution
			json.Unmarshal(v, &e)
			executions = append(executions, e)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (b *Bolt) UpdateCandle(c models.Candle) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucketName := getCandleBucketName(c.Duration)
		durationBucket, err := tx.CreateBucketIfNotExists([]byte(bucketName ))
		if err != nil {
			return err
		}
		bucket,err := durationBucket.CreateBucketIfNotExists([]byte(candleBucket))
		if err != nil {
			return err
		}
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

func (b *Bolt) GetCandle(duration string) models.Candle {
	return models.Candle{}
}

func (b *Bolt) GetCandleCollection() {

}
