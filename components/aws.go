package components

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitAws initializes and caches the aws session
// reads the AWS region from the environment variable AWS_REGION, example: eu-west-1
func InitAws() (err error) {
	var sess *session.Session
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return err
	}

	cache.SetAwsSession(sess)
	return nil
}
