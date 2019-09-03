package db

import (
	"github.com/boltdb/bolt"
	"encoding/json"
	"github.com/atoyr/goflyer/models"
"bytes"
"fmt"
"time"
"log"
)

type Bolt struct {
	db *bolt.DB
}

const TickerBucket = "Ticker"

func GetBolt(dbFile string) (*Bolt, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	bolt := new(Bolt)
	bolt.db = db
	return bolt, nil
}

func (b *Bolt) Init() error {
	b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TickerBucket))
		_, err = tx.CreateBucketIfNotExists([]byte("Candle"))
		return err
	})
	return nil
}

func (b *Bolt) UpdateTicker(t models.Ticker) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket  := tx.Bucket([]byte(TickerBucket))
		if marshalTime, err := t.GetTimestamp().MarshalBinary() ; err!= nil {
			log.Fatal("hoge")
			return err 
		} else if buf, err := json.Marshal(t); err != nil {
			log.Fatal("fugu")
			return err
		} else if err = bucket.Put(marshalTime,buf) ; err != nil {
			log.Fatal("piyo")
			return err
		}
		return nil
	})
		if err  != nil {
			log.Fatal(err)
			return err
		}
		return nil
}

func (b *Bolt) GetTicker(timestamp time.Time) (*models.Ticker, error) {
	m := new(models.Ticker)
	err := b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(TickerBucket)).Cursor()
		min ,err  :=  timestamp.MarshalBinary()
		if err  != nil {
			return err
		}
		max ,err :=  timestamp.MarshalBinary()
		if err  != nil {
			return err
		}

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v) 
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return m, nil
}
