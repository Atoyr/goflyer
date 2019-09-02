package db

import(
	"github.com/boltdb/bolt"
)

type Bolt struct {
	db *bolt.DB

}


func GetBolt(dbFile string )  (*Bolt, error){ 
	db, err := bolt.Open(dbFile, 0600, nil )
	if err != nil {
		return nil ,err
	}
	bolt := new(Bolt)
	bolt.db = db
	return bolt,nil
}

func (b *Bolt) Init() error {
	b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Ticker"))
		return nil
	})
	return nil
}
