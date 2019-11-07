package executor

import (
	"sync"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/configs"
)

// Executor is singleton
type executor struct {
	dataFrames []models.DataFrame
	db         db.DB
	client     client.APIClient
}

var (
	once sync.Once
	exe  *executor
)

func getExecutor() *executor {
	once.Do(func() {
		config ,err := configs.GetGeneralConfig()
		if err != nil {
			panic(err)
		}
		e := new(executor)
		e.dataFrames = make([]models.DataFrame, 0)
		e.client = *client.New(config.Apikey(),config.Secretkey())

		e.db = config.GetDB()

		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("1m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("3m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("5m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("10m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("15m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("30m")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("1h")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("2h")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("4h")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("6h")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("12h")))
		e.dataFrames = append(e.dataFrames, e.db.GetDataFrame(models.GetDuration("24h")))
		exe = e
	})
	return exe
}

func ChangeDB(db db.DB) {
	exe := getExecutor()
	exe.db = db
}

