package cache

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	awsSession       *session.Session
	awsSqsSession    *sqs.SQS
	awsLambdaSession *lambda.Lambda
	awsSessionMutex  sync.RWMutex
)

// SetAwsSession caches an aws session for future use
func SetAwsSession(s *session.Session) {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	awsSession = s
}

// GetAwsSession returns a cached aws session
func GetAwsSession() *session.Session {
	awsSessionMutex.RLock()
	defer awsSessionMutex.RUnlock()

	return awsSession
}

// SetAwsSqsSession caches an aws sqs session for future use
func SetAwsSqsSession(s *sqs.SQS) {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	awsSqsSession = s
}

// GetAwsSqsSession returns a cached aws sqs session
func GetAwsSqsSession() *sqs.SQS {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	return awsSqsSession
}

// SetAwsLambdaSession caches an aws lambda session for future use
func SetAwsLambdaSession(s *lambda.Lambda) {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	awsLambdaSession = s
}

// GetAwsLambdaSession returns a cached aws lambda session
func GetAwsLambdaSession() *lambda.Lambda {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	return awsLambdaSession
}
