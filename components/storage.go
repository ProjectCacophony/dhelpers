package components

import (
	"os"

	"github.com/minio/minio-go"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitMinio sets up and caches the minio client
// reads the s3 endpoint from S3_ENDPOINT
// reads the s3 access key from S3_ACCESS_KEY
// reads the s3 access secret from S3_ACCESS_SECRET
// reads the s3 bucket from S3_BUCKET
// reads the s3 location from S3_LOCATION
// reads the s3 cache folder from S3_CACHE_FOLDER
// if S3_NOTSECURE is set, it will connect to the s3 server insecure
func InitMinio() (err error) {
	var minioClient *minio.Client

	_, notsecureExists := os.LookupEnv("S3_NOTSECURE")
	secure := !notsecureExists

	minioClient, err = minio.New(
		os.Getenv("S3_ENDPOINT"),
		os.Getenv("S3_ACCESS_KEY"),
		os.Getenv("S3_ACCESS_SECRET"),
		secure,
	)
	if err != nil {
		return err
	}

	cache.SetMinio(minioClient)

	bucketExists, err := minioClient.BucketExists(os.Getenv("S3_BUCKET"))
	if err != nil {
		return err
	}

	if !bucketExists {
		err = minioClient.MakeBucket(os.Getenv("S3_BUCKET"), os.Getenv("S3_LOCATION"))
		if err != nil {
			return err
		}
	}

	return err
}
