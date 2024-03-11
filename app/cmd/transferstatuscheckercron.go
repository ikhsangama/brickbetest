package cmd

import (
	"brickbetest/config"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

func StartTransferStatusCheckerCron() {
	c := cron.New(cron.WithSeconds())
	cronCfg := config.GetTransferStatusCheckerConfig()
	log.Println(cronCfg)
	app := initApp()
	ctx := context.Background()
	entryId, err := c.AddFunc(cronCfg.CronSpec, func() {
		log.Println("Running bank transfer status checker cron..")
		err := app.transferService.TransferStatusCheck(ctx, cronCfg.IntervalDays, cronCfg.FetchLimit)
		if err != nil {
			fmt.Printf("failed to check transfer status %v", err)
		}
		log.Println("Transfer status checker finished.")
	})
	if err != nil {
		return
	}

	fmt.Print(entryId)

	c.Start()

	// This will block indefinitely until the program is explicitly killed
	terminate := make(chan struct{})
	<-terminate
}
