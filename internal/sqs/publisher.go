package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
)

type Publisher struct {
	client *sqs.Client
	url    string
}

func NewSqsPublisher(url string) *Publisher {
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
	return &Publisher{
		client: client,
		url:    url,
	}
}

func (p *Publisher) Publish(ctx context.Context, message interface{}, path string) error {
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal %s", err)
		return err
	}
	url := fmt.Sprintf(p.url + path)
	messageString := string(messageByte)
	log.Printf(messageString)

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  &messageString,
		QueueUrl:     &url,
		DelaySeconds: 2,
	})
	if err != nil {
		log.Printf("Failed to send message: %v. Error: %v", messageString, err)
		return err
	}
	return nil
}
