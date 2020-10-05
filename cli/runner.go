package cli

import (
	"context"
  "fmt"
  "bufio"
  "os"
  "time"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/controllers"
	urfavecli "github.com/urfave/cli"
)

func runCommand() *urfavecli.Command {
	var command urfavecli.Command
	command.Name = "run"
	command.Aliases = []string{"r"}
	command.Action = runAction

	return &command
}

func runAction(clictx *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

  c := client.New("","")
  cc := controllers.NewClientController(*c)

	cc.SubscribeTicker(func(ticker bitflyer.Ticker) {
		fmt.Printf("\r%s value %f ", ticker.Timestamp, ticker.Ltp)
	})

  go cc.ExecuteFetchTicker(ctx)

  for {
    fmt.Println("application end is typed end")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    s := scanner.Text()
    switch s {

    case "start":
      cc.RegisterScheduleAction(time.Now(),controllers.StartAction)
      break;
    case "stop":
      cc.RegisterScheduleAction(time.Now(),controllers.StopAction)
      break;
    case "end":
      cc.RegisterScheduleAction(time.Now(),controllers.ExitAction)
      break;
    }
  }

	return nil
}
