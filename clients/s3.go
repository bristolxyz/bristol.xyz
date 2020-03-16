package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bristolxyz/bristol.xyz/utils"
)

// CDNHostname defines the CDN hostname.
var CDNHostname string

// S3Instance is the exported S3 instance.
var S3Instance *s3.S3

// S3Init is used to initialse S3.
func S3Init() {
	EnvMap := utils.RequiredEnvs(
		"S3_ENDPOINT", "S3_REGION", "S3_BUCKET", "CDN_HOSTNAME", "AWS_SECRET_ACCESS_KEY",
		"AWS_ACCESS_KEY_ID")
	CDNHostname = EnvMap["CDN_HOSTNAME"]
	StaticCredential := credentials.NewStaticCredentials(EnvMap["AWS_ACCESS_KEY_ID"], EnvMap["AWS_SECRET_ACCESS_KEY"], "")
	e := EnvMap["S3_ENDPOINT"]
	r := EnvMap["S3_REGION"]
	s3sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &e,
		Credentials: StaticCredential,
		Region:      &r,
	}))
	S3Instance = s3.New(s3sess)
}
