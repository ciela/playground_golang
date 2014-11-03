package aws

import (
	"log"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var LgtmBucket *s3.Bucket

func init() {
	auth, err := aws.EnvAuth() //TODO setenv auths
	if err != nil {
		log.Fatalln(err)
	}
	s3 := s3.New(auth, aws.APNortheast)
	LgtmBucket = s3.Bucket("lgtm-images")
}
