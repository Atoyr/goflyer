package cli

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/models"
	urfavecli "github.com/urfave/cli"
)

func runWebappsCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "runweb"
	command.Action = runWebappsAction

	return command
}

func runWebappsAction(c *urfavecli.Context) error {
	jsonFile, err := os.Open("./testdata/tickers.json")
	if err != nil {
		log.Println(err)
		return err
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	cc := models.NewDataFrame(models.BTC_JPY,models.GetDuration("3m"))
		tickers, err := models.JsonUnmarshalTickers(raw)
	if err != nil {
		return err
	}
	for i := range tickers {
		cc.AddTicker(tickers[i])
	}
	cc.AddEmas(6)
	ccs := models.DataFrames{}
	ccs[cc.Name()] = cc
	e := api.AppendHandler(api.GetEcho(ccs))
	e.Start(":8080")
	return nil
}
