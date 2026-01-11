package mb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"path"
	"strings"
)

func (request *MultibeamRequest) resolveSurveys() {
	fmt.Println("Resolving bathymetry data for specified surveys: ", request.Surveys)
	var surveyPaths []string
	wantedSurveys := len(request.Surveys)
	foundSurveys := 0

	pt, ptErr := request.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(multibeamNODDBucket),
		Prefix:    aws.String("mb/"),
		Delimiter: aws.String("/"),
	})

	if ptErr != nil {
		request.Error = ptErr
		return
	}

	for _, platformType := range pt.CommonPrefixes {

		platformParams := &s3.ListObjectsV2Input{
			Bucket:    aws.String(multibeamNODDBucket),
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
					Bucket:    aws.String(multibeamNODDBucket),
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
