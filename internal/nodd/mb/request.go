package mb

import "github.com/aws/aws-sdk-go-v2/service/s3"

type MultibeamRequest struct {
	Surveys     []string
	Prefixes    []string
	S3Client    s3.Client
	TargetDir   string
	WorkerCount int
	Error       error
}
