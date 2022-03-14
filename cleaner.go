package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage : sqs-cleaner https://sqs.<REGION>.amazonaws.com/<ACCOUNT-ID>/<NAME>")
		os.Exit(1)
	}

	log.Printf("Queue URL : %s", os.Args[1])
	queueUrl := aws.String(os.Args[1])

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		os.Exit(1)
	}

	svc := sqs.NewFromConfig(cfg)

	for {
		log.Printf("ReceiveMessage...")
		res, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: queueUrl,
		})

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		for _, msg := range res.Messages {
			log.Printf("[%s] %s", *msg.MessageId, *msg.Body)
			svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      queueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
	log.Printf("Terminate")
}
