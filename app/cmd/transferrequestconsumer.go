package cmd

import (
	"brickbetest/config"
	"brickbetest/internal/sqs"
	"brickbetest/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

func StartTransferRequestConsumer() {
	app := initApp()

	consumer := sqs.NewSQSConsumer(config.GetSqsUrl() + sqs.BankTransferRequest)
	ctx := context.Background()

	err := consumer.Handle(ctx, func(message types.Message) sqs.HandlerOutput {
		var txn model.Transfer
		err := json.Unmarshal([]byte(*message.Body), &txn)
		if err != nil {
			log.Printf("failed to process message %v , error: %v", message, err)
			return sqs.HandlerOutput{}
		}

		log.Printf("handling %v", txn)
		err = app.transferService.HandleTransferRequest(ctx, txn)
		if err != nil {
			log.Printf("failed handling message %v", err)
		}

		return sqs.HandlerOutput{
			Ack: err == nil,
		}
	})

	if err != nil {
		panic(fmt.Sprintf("failed to start consumer %v", err))
	}
}
