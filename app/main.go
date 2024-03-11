package main

import (
	cmd2 "brickbetest/app/cmd"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected one of these subcommands: 'start-rest-server', 'start-transferrequest-consumer', 'start-recordtxn-consumer', 'start-transferstatuschecker-cron'")
		os.Exit(1)
	}

	var cmd = os.Args[1]

	if cmd == "start-rest-server" {
		cmd2.StartRestServer()
	} else if cmd == "start-transferrequest-consumer" {
		cmd2.StartTransferRequestConsumer()
	} else if cmd == "start-recordtxn-consumer" {
		cmd2.StartRecordTxnConsumer()
	} else if cmd == "start-transferstatuschecker-cron" {
		cmd2.StartTransferStatusCheckerCron()
	} else {
		fmt.Println("Invalid subcommand. Expected one of these subcommands: 'start-rest-server', 'start-transferrequest-consumer', 'start-recordtxn-consumer', 'start-transferstatuschecker-cron'")
		os.Exit(1)
	}
}
