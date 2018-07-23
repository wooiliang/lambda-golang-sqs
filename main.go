package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/teris-io/shortid"
	"github.com/wooiliang/aws-lambda-go/events"
	"github.com/wooiliang/aws-lambda-go/lambda"
)

// HelloWorld struct
type HelloWorld struct {
	Foo string `json:"foo"`
}

func putItem(helloWorld *HelloWorld) error {
	shortID, _ := shortid.Generate()
	svc := dynamodb.New(session.New())
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(shortID),
			},
			"foo": {
				S: aws.String(helloWorld.Foo),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("EmployerActivities"),
	}
	result, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
				return aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}
	fmt.Println(result)
	return nil
}

func getJSON(body string) (*HelloWorld, error) {
	helloWorld := &HelloWorld{}
	if err := json.Unmarshal([]byte(body), helloWorld); err != nil {
		return nil, err
	}
	return helloWorld, nil
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		if helloWorld, err := getJSON(message.Body); err != nil {
			fmt.Println(err)
			return err
		} else if err := putItem(helloWorld); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
