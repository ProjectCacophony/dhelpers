package dhelpers

import (
	"encoding/binary"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/json-iterator/go"
)

// StartLambdaAsync sends an event to the given AWS Lambda Function
func StartLambdaAsync(lambdaClient *lambda.Lambda, event EventContainer, function string) (bytesSent int, err error) {
	// pack the event data
	marshalled, err := jsoniter.Marshal(event)
	if err != nil {
		return 0, err
	}

	// invoke lambda
	_, err = lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(function),
		InvocationType: aws.String("Event"), // Async
		Payload:        marshalled,
	})
	if err != nil {
		return 0, errors.New("error invoking lambda: " + err.Error())
	}
	return binary.Size(marshalled), nil
}
