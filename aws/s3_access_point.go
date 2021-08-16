package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/gruntwork-io/cloud-nuke/config"
)

func getAllBucketAccessPoints(awsSession *session.Session, bucketName string, cfgObj *config.Config) (accesspoints *s3control.ListAccessPointsOutput, err error) {
	svc := s3control.New(awsSession)

	return svc.ListAccessPoints(&s3control.ListAccessPointsInput{
		AccountId: aws.String(cfgObj.AWSAccountID),
		Bucket:    aws.String(bucketName),
	})

}

func deleteAccessPoint(awsSession *session.Session, AccessPointName *string, cfgObj *config.Config) (*s3control.DeleteAccessPointOutput, error) {
	svc := s3control.New(awsSession)
	return svc.DeleteAccessPoint(&s3control.DeleteAccessPointInput{
		AccountId: aws.String(cfgObj.AWSAccountID),
		Name:      AccessPointName,
	})
}
