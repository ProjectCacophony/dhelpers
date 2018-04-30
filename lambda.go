package dhelpers

import (
	"context"
	"encoding/binary"
	"errors"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/json-iterator/go"
	"gitlab.com/Cacophony/dhelpers/cache"
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

type lambdaHandler func(context.Context, []byte) ([]byte, error)

func (handler lambdaHandler) Invoke(_ context.Context, payload []byte) ([]byte, error) {
	_, err := handler(nil, payload)
	return nil, err
}

// NewLambdaHandler initialises a lambda event receiver
func NewLambdaHandler(service string, method func(event EventContainer)) lambdaHandler { // nolint: golint
	return func(_ context.Context, payload []byte) ([]byte, error) {
		// unmarshal
		unmarshalStart := time.Now()
		var container EventContainer
		err := jsoniter.Unmarshal(payload, &container)
		if err != nil {
			cache.GetLogger().Errorln("error unmarshalling event container", err.Error())
			return nil, err
		}
		cache.GetLogger().Infoln("unmarshal took", time.Since(unmarshalStart).String())

		// check if event is valid
		if len(container.Destinations) < 1 {
			cache.GetLogger().Errorln("invalid event received, no destinations")
			return nil, errors.New("invalid event received, no destinations")
		}

		// error handling
		defer func() {
			err := recover()
			if err != nil {
				HandleErrWith(service, err.(error), container.Destinations[0].ErrorHandlers, &container)
			}
		}()

		// start handler
		method(container)
		return nil, nil
	}
}
