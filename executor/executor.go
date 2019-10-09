package executor

import (
	"sync"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
)

type executor struct {
	dataFrames models.DataFrames
	db         db.DB
}

var (
	once sync.Once
	exe  *executor
)

func GetExecutor(db db.DB) *executor {
	once.Do(func() {
		e := new(executor)
		e.dataFrames = make(map[string]models.DataFrame, 0)
		e.db = db
		exe = e
	})
	return exe
}
