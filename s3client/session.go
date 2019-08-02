package s3client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(s3Config S3Config) *session.Session {
	config := &aws.Config{
		Region:           aws.String(s3Config.Region),
		Endpoint:         aws.String(s3Config.EndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(s3Config.ID, s3Config.Secret, ""),
	}
	return session.Must(session.NewSession(config))
}
