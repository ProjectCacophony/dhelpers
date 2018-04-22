package components

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitAwsLambda initializes and caches an aws lambda session
// will initialize AWS if AWS has not been initialized yet
// reads the AWS region from the environment variable AWS_REGION, example: eu-west-1
func InitAwsLambda() (err error) {
	if cache.GetAwsSession() == nil {
		err = InitAws()
		if err != nil {
			return err
		}
	}

	lambdaClient := lambda.New(cache.GetAwsSession())

	cache.SetAwsLambdaSession(lambdaClient)

	return nil
}
