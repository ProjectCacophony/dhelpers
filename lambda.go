package dhelpers

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/json-iterator/go"
)

// sends an event to the given AWS Lambda Function
func StartLambdaAsync(lambdaClient *lambda.Lambda, start, receive time.Time, eventType EventType, eventData interface{}, function string) (bytesSent int, err error) {
	// pack the event data
	marshalled, err := jsoniter.Marshal(eventData)
	if err != nil {
		return 0, err
	}

	// create event container
	eventContainer := EventContainer{
		Type:           eventType,
		ReceivedAt:     receive,
		GatewayStarted: start,
		Data:           marshalled,
	}
	// pack the event container
	marshalledContainer, err := jsoniter.Marshal(eventContainer)
	if err != nil {
		return 0, err
	}

	// invoke lambda
	_, err = lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(function),
		InvocationType: aws.String("Event"), // Async
		Payload:        marshalledContainer,
	})
	if err != nil {
		return 0, errors.New("error invoking lambda: " + err.Error())
	}
	return binary.Size(marshalledContainer), nil
}
