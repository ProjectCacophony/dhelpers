package components

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitAwsSqs initializes and caches an aws sqs session
// will initialize AWS if AWS has not been initialized yet
// reads the AWS region from the environment variable AWS_REGION, example: eu-west-1
func InitAwsSqs() (err error) {
	if cache.GetAwsSession() == nil {
		err = InitAws()
		if err != nil {
			return err
		}
	}

	sqsClient := sqs.New(cache.GetAwsSession())

	cache.SetAwsSqsSession(sqsClient)

	return nil
}
