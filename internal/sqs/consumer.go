package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
	"time"
)

type HandlerOutput struct {
	Ack bool
}

type Consumer struct {
	client *sqs.Client
	url    string
}

func NewSQSConsumer(url string) *Consumer {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	client := sqs.NewFromConfig(cfg)
	return &Consumer{
		url:    url,
		client: client,
	}
}

func (p *Consumer) Handle(ctx context.Context, handler func(message types.Message) HandlerOutput) error {
	var err error
	var receiveOutput *sqs.ReceiveMessageOutput
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			receiveOutput, err = p.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:        &p.url,
				WaitTimeSeconds: 10,
			})
			if err != nil {
				log.Printf("failed to received a message from SQS: %v", err)
				time.Sleep(1 * time.Second) // To avoid overloading the Retry immediately
				continue
			}

			for _, message := range receiveOutput.Messages {
				output := handler(message)
				if output.Ack {
					_, err = p.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      &p.url,
						ReceiptHandle: message.ReceiptHandle,
					})
					if err != nil {
						log.Printf("failed to delete message from queue: %v", err)
						time.Sleep(1 * time.Second) // To avoid overloading the Retry immediately
						continue
					}
					log.Printf("message processed")
				}
			}
		}
	}
}
