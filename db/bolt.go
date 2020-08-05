package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/models/bitflyer"
	"github.com/atoyr/goflyer/util"
	"github.com/boltdb/bolt"
)

type Bolt struct {
	dbFile string
}

// Bucket layout
// - tickerBucket
// - executionBuclet
// - N - DurationBucket
//   - OpenBucket
//   - CloseBucket
//   - HighBucket
//   - LowBucket

const (
	tickerBucketName    = "Ticker"
	executionBucketName = "Execution"
	durationBucketName  = "Duration"
	logBucketName       = "Log"
)

func GetBolt(dbFile string) (Bolt, error) {
	var b Bolt
	b.dbFile = dbFile
	err := b.init()
	return b, err
}

func getCandleBucketName(duration time.Duration) string {
	return fmt.Sprintf("%s_%s", durationBucketName, models.GetDurationString(duration))
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

func (b *Bolt) UpdateTicker(t bitflyer.Ticker) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tickerBucketName))
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

func (b *Bolt) GetTicker(tickID float64) (bitflyer.Ticker, error) {
	db := b.db()
	defer db.Close()
	ticker := new(bitflyer.Ticker)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(tickerBucketName)).Cursor()

		marshalID := util.Float64ToBytes(tickID)

		for k, v := c.Seek(marshalID); k != nil && bytes.Compare(k, marshalID) <= 0; k, v = c.Next() {
			t, err := bitflyer.JsonUnmarshalTicker(v)
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

func (b *Bolt) GetTickerAll() ([]bitflyer.Ticker, error) {
	db := b.db()
	defer db.Close()
	tickers := make([]bitflyer.Ticker, 0)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(tickerBucketName)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t bitflyer.Ticker
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

func (b *Bolt) UpdateExecution(execution bitflyer.Execution) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(executionBucketName))
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
func (b *Bolt) GetExecutionAll() ([]bitflyer.Execution, error) {
	db := b.db()
	defer db.Close()
	executions := make([]bitflyer.Execution, 0)
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(executionBucketName)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e bitflyer.Execution
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
func (b *Bolt) GetCandles(duration time.Duration) (models.Candles, error) {
	db := b.db()
	defer db.Close()
	cs := models.NewCandles("BTC_JPY", time.Duration(duration))
	err := db.View(func(tx *bolt.Tx) error {
		bucketName := getCandleBucketName(duration)
		durationBucket := tx.Bucket([]byte(bucketName))
		if durationBucket == nil {
			return fmt.Errorf("%s bucket not found", bucketName)
		}
		bucket := durationBucket.Bucket([]byte("candleBucket"))
		if bucket == nil {
			return fmt.Errorf("candle bucket not found")
		}
		err := bucket.ForEach(func(k, v []byte) error {
			_, err := models.JsonUnmarshalCandle(v)
			if err != nil {
				return err
			}
			// TODO append candle for candles
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return cs, err
	}
	return cs, nil
}

func (b *Bolt) UpdateCandle(duration time.Duration, c models.Candle) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucketName := getCandleBucketName(duration)
		durationBucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		bucket, err := durationBucket.CreateBucketIfNotExists([]byte("candleBucket"))
		if err != nil {
			return err
		}
		if buf, err := json.Marshal(c); err != nil {
			return err
		} else if err = bucket.Put([]byte(fmt.Sprintf("%v", c)), buf); err != nil {
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

func (b *Bolt) UpdateDataFrame(df models.DataFrame) error {
	if len(df.Datetimes) == 0 {
		return nil
	}
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		durationBucket, err := tx.CreateBucketIfNotExists([]byte(models.GetDurationString(df.Duration)))
		if err != nil {
			return err
		}
		fromIndex := -1
		lastTime := durationBucket.Get([]byte("tail"))
		if lastTime != nil {
			var t time.Time
			err = t.UnmarshalBinary(lastTime)
			if err != nil {
				return err
			}
			for i := range df.Datetimes {
				// TODO bug it !!!
				if t.Before(df.Datetimes[i]) {
					fromIndex = i
					break
				}
				if t.Equal(df.Datetimes[i]) {
					fromIndex = i
					break
				}
			}
			if fromIndex < 0 {
				fmt.Println(t)
				return nil
			}
		} else {
			fromIndex = 0
		}
		openBucket, err := durationBucket.CreateBucketIfNotExists([]byte("open"))
		if err != nil {
			return err
		}
		closeBucket, err := durationBucket.CreateBucketIfNotExists([]byte("close"))
		if err != nil {
			return err
		}
		highBucket, err := durationBucket.CreateBucketIfNotExists([]byte("high"))
		if err != nil {
			return err
		}
		lowBucket, err := durationBucket.CreateBucketIfNotExists([]byte("low"))
		if err != nil {
			return err
		}
		volumeBucket, err := durationBucket.CreateBucketIfNotExists([]byte("volume"))
		if err != nil {
			return err
		}

		for i := range df.Datetimes[fromIndex:] {
			index := i + fromIndex
			key, err := df.Datetimes[index].MarshalBinary()
			if err != nil {
				return err
			}
			open := util.Float64ToBytes(df.Opens[index])
			close := util.Float64ToBytes(df.Closes[index])
			high := util.Float64ToBytes(df.Highs[index])
			low := util.Float64ToBytes(df.Lows[index])
			volume := util.Float64ToBytes(df.Volumes[index])
			err = openBucket.Put(key, []byte(open))
			if err != nil {
				return err
			}
			err = closeBucket.Put(key, []byte(close))
			if err != nil {
				return err
			}
			err = highBucket.Put(key, []byte(high))
			if err != nil {
				return err
			}
			err = lowBucket.Put(key, []byte(low))
			if err != nil {
				return err
			}
			err = volumeBucket.Put(key, []byte(volume))
			if err != nil {
				return err
			}
		}
		tailTime, err := df.Datetimes[len(df.Datetimes)-1].MarshalBinary()
		if err != nil {
			return err
		}
		durationBucket.Put([]byte("tail"), tailTime)
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (b *Bolt) GetDataFrame(duration time.Duration) models.DataFrame {
	df := models.NewDataFrame("BTC_JPY", duration)
	db := b.db()
	defer db.Close()
	err := db.View(func(tx *bolt.Tx) error {
		durationBucket := tx.Bucket([]byte(models.GetDurationString(df.Duration)))
		if durationBucket == nil {
			return errors.New("duration bucket not found")
		}
		openBucket := durationBucket.Bucket([]byte("open"))
		if openBucket == nil {
			return errors.New("open bucket not found")
		}
		closeBucket := durationBucket.Bucket([]byte("close"))
		if closeBucket == nil {
			return errors.New("close bucket not found")
		}
		highBucket := durationBucket.Bucket([]byte("high"))
		if highBucket == nil {
			return errors.New("high bucket not found")
		}
		lowBucket := durationBucket.Bucket([]byte("low"))
		if lowBucket == nil {
			return errors.New("low bucket not found")
		}
		volumeBucket := durationBucket.Bucket([]byte("volume"))
		if volumeBucket == nil {
			return errors.New("volume bucket not found")
		}
		datetimes := make([]time.Time, 0)
		opens := make([]float64, 0)
		closes := make([]float64, 0)
		highs := make([]float64, 0)
		lows := make([]float64, 0)
		volumes := make([]float64, 0)
		openBucket.ForEach(func(k, v []byte) error {
			var t time.Time
			err := t.UnmarshalBinary(k)
			if err != nil {
				return err
			}
			datetimes = append(datetimes, t)
			o := util.BytesToFloat64(v)
			opens = append(opens, o)
			return nil
		})
		closeBucket.ForEach(func(k, v []byte) error {
			c := util.BytesToFloat64(v)
			closes = append(closes, c)
			return nil
		})
		highBucket.ForEach(func(k, v []byte) error {
			h := util.BytesToFloat64(v)
			highs = append(highs, h)
			return nil
		})
		lowBucket.ForEach(func(k, v []byte) error {
			l := util.BytesToFloat64(v)
			lows = append(lows, l)
			return nil
		})
		volumeBucket.ForEach(func(k, v []byte) error {
			vol := util.BytesToFloat64(v)
			volumes = append(volumes, vol)
			return nil
		})
		df.Datetimes = datetimes
		df.Opens = opens
		df.Closes = closes
		df.Highs = highs
		df.Lows = lows
		df.Volumes = volumes
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return df
}

func (b *Bolt) AddLog(datetime time.Time, log models.Log) error {
	db := b.db()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		logBucket, err := tx.CreateBucketIfNotExists([]byte(logBucketName))
		if err != nil {
			return err
		}
		key, err := datetime.MarshalBinary()
		if err != nil {
			return err
		}
		marshalLog, err := json.Marshal(log)
		if err != nil {
			return err
		}
		logBucket.Put(key, marshalLog)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Bolt) GetLog(from, to time.Time) ([]models.Log, error) {
	db := b.db()
	defer db.Close()
	logs := make([]models.Log, 0)
	err := db.View(func(tx *bolt.Tx) error {
		logBucket := tx.Bucket([]byte(logBucketName))
		if logBucket == nil {
			return fmt.Errorf("log Bucket not found")
		}
		f, err := from.MarshalBinary()
		if err != nil {
			return err
		}
		t, err := to.MarshalBinary()
		if err != nil {
			return err
		}
		c := logBucket.Cursor()
		for k, v := c.Seek(f); k != nil && bytes.Compare(k, t) <= 0; k, v = c.Next() {
			log, err := models.JsonUnmarshalLog(v)
			if err != nil {
				return err
			}
			logs = append(logs, *log)
		}
		return nil
	})
	if err != nil {
		return logs, err
	}
	return logs, nil

}
