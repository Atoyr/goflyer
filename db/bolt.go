package db

import (
	"github.com/boltdb/bolt"
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
		b, err := tx.CreateBucketIfNotExists([]byte(TickerBucket))
		b, err := tx.CreateBucketIfNotExists([]byte("Candle"))
		return nil
	})
	return nil
}

func (b *Bolt) UpdateTicker(t model.Ticker) {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.Bucket([]byte(TickerBucket))
		if buf, err := json.Marshal(t); err != nil {
			return err
		} else if err := bucket.Put([]byte(t.GetTimestamp()), buf); err != nil {
			return err
		}
	})
}

func (b *Bolt) GetTicker(timestamp time.Time) (*model.Ticker, error) {
	m := new(model.Ticker)
	err := b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(TickerBucket)).Cursor()
		min := []byte(timestamp)
		max := []byte(timestamp)

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
