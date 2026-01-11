package mb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/max-e-smith/cruise-lug/internal/common"
	"log"
	"path"
	"strings"
	"time"
)

var BathyBucket = "noaa-dcdb-bathymetry-pds" // https://noaa-dcdb-bathymetry-pds.s3.amazonaws.com/index.html

type MultibeamRequest struct {
	Surveys     []string
	Prefixes    []string
	S3Client    s3.Client
	TargetDir   string
	WorkerCount int
	Error       error
}

func RequestMultibeamDownload(request MultibeamRequest) {
	request.resolveSurveys()
	request.checkDiskAvailability()
	request.downloadSurveys()

	if request.Error != nil {
		log.Fatal(request.Error)
	}
}

func logDownloadTime(start time.Time) {
	fmt.Printf("Download completed in %g hours.\n", common.HoursSince(start))
}

func (request *MultibeamRequest) resolveSurveys() {
	fmt.Println("Resolving bathymetry data for specified surveys: ", request.Surveys)
	var surveyPaths []string
	wantedSurveys := len(request.Surveys)
	foundSurveys := 0

	pt, ptErr := request.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(BathyBucket),
		Prefix:    aws.String("mb/"),
		Delimiter: aws.String("/"),
	})

	if ptErr != nil {
		request.Error = ptErr
		return
	}

	for _, platformType := range pt.CommonPrefixes {

		platformParams := &s3.ListObjectsV2Input{
			Bucket:    aws.String(BathyBucket),
			Prefix:    aws.String(*platformType.Prefix),
			Delimiter: aws.String("/"),
		}

		allPlatforms := s3.NewListObjectsV2Paginator(&request.S3Client, platformParams)

		for allPlatforms.HasMorePages() {
			platsPage, platsErr := allPlatforms.NextPage(context.TODO())

			if platsErr != nil {
				request.Error = platsErr
				return
			}

			for _, platform := range platsPage.CommonPrefixes {
				fmt.Printf("  searching %s\n", *platform.Prefix)

				platformParams := &s3.ListObjectsV2Input{
					Bucket:    aws.String(BathyBucket),
					Prefix:    aws.String(*platform.Prefix),
					Delimiter: aws.String("/"),
				}

				platformPaginator := s3.NewListObjectsV2Paginator(&request.S3Client, platformParams)

				for platformPaginator.HasMorePages() {
					surveysPage, err := platformPaginator.NextPage(context.TODO())
					if err != nil {
						request.Error = err
						return
					}

					for _, survey := range surveysPage.CommonPrefixes {
						surveyPrefix := *survey.Prefix
						survey := path.Base(strings.TrimRight(surveyPrefix, "/"))
						if isSurveyMatch(request.Surveys, survey) {
							surveyPaths = append(surveyPaths, surveyPrefix)
							foundSurveys++
						}
					}

				}

				if wantedSurveys == foundSurveys {
					// all surveys are found
					request.Prefixes = surveyPaths
					return
				}
			}
		}
	}

	if len(surveyPaths) == 0 {
		fmt.Printf("No matching surveys found for %s\n", request.Surveys)
	} else {
		fmt.Printf("Found %d of %d wanted surveys at: %s\n", len(surveyPaths), len(request.Surveys), surveyPaths)
		request.Prefixes = surveyPaths

	}
	return
}

func isSurveyMatch(surveys []string, resolvedSurvey string) bool {
	for _, survey := range surveys {
		if survey == resolvedSurvey {
			fmt.Println("Found matching survey: ", survey)
			return true
		}
	}
	return false
}

func (request *MultibeamRequest) checkDiskAvailability() {
	if request.Error != nil || len(request.Prefixes) == 0 {
		return
	}
	// TODO get viper check config and return if false

	bytes, estimateErr := common.GetDiskUsageEstimate(BathyBucket, request.S3Client, request.Prefixes)
	if estimateErr != nil {
		request.Error = errors.Join(errors.New("unable to get disk usage estimate from s3 bucket"), estimateErr)
		return
	}

	spaceErr := common.DiskSpaceCheck(bytes, request.TargetDir)
	if spaceErr != nil {
		request.Error = spaceErr
	}

	return
}

func (request *MultibeamRequest) downloadSurveys() {
	if request.Error != nil || len(request.Prefixes) == 0 {
		return
	}

	start := time.Now()
	defer logDownloadTime(start)

	order := common.Order{
		Bucket:      BathyBucket,
		Prefixes:    request.Prefixes,
		Client:      request.S3Client,
		TargetDir:   request.TargetDir,
		WorkerCount: request.WorkerCount,
	}

	if err := order.DownloadFiles(); err != nil {
		request.Error = err
	}

	return
}
