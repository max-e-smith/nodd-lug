package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"log"
)

var bathy bool
var wcd bool
var trackline bool

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download NOAA survey data to local path",
	Long: `Use 'clug get <survey(s)> <local path> <options>' to download marine geophysics data to your machine. 

		Data is downloaded from the NOAA Open Data Dissemination cloud buckets by default. You must 
		specify a data type(s) for this command. View the help for more info on those options. Specify
		the survey(s) you want to download and a local path to download data to. The path must exist and 
		have the necessary permissions.`,
	Run: func(cmd *cobra.Command, args []string) {

		var length = len(args)
		if length <= 1 {
			fmt.Println("Please specify survey name(s) and a target file path.")
			fmt.Println(cmd.UsageString())
			return
		}

		var path = args[length-1]
		var surveys = args[:length-1] // high is non-inclusive

		if !bathy && !wcd && !trackline {
			fmt.Println("Please specify data type(s) for download.")
			fmt.Println(cmd.UsageString())
			return
		}

		download(surveys, path)

		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Local flags
	getCmd.Flags().BoolVarP(&bathy, "bathy", "b", false, "Download bathy data")
	getCmd.Flags().BoolVarP(&wcd, "water-column", "w", false, "Download water column data")
	getCmd.Flags().BoolVarP(&trackline, "trackline", "t", false, "Download trackline data")

}

func download(surveys []string, path string) {

	if bathy {
		fmt.Println("resolving bathymetry data for provided surveys: ", surveys)
		var surveyRoots []string = resolveBathySurveys(surveys)

		if len(surveyRoots) == 0 {
			fmt.Println("No surveys found.")
			return
		}

		fmt.Println("checking available disk space")
		diskSpaceCheck(surveyRoots)

		fmt.Println("downloading surveys ", surveys, " to ", path, "...")
		// TODO recursively download surveys

		fmt.Println("bathymetry data downloaded.")
	}

	if wcd {
		// TODO
	}

	if trackline {
		// TODO
	}

	return
}

func diskSpaceCheck(rootPaths []string) {
	// TODO
}

func resolveBathySurveys(surveys []string) []string {
	var surveyPaths []string

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
		return surveyPaths
	}

	client := s3.NewFromConfig(cfg)

	// noaa-dcdb-bathymetry-pds.s3.amazonaws.com/index.html
	bucket := "noaa-dcdb-bathymetry-pds"

	pt, ptErr := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String("mb/"),
		Delimiter: aws.String("/"),
	})

	if ptErr != nil {
		log.Fatal(err)
		return surveyPaths
	}

	for _, platformType := range pt.CommonPrefixes {

		platformParams := &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucket),
			Prefix:    aws.String(*platformType.Prefix),
			Delimiter: aws.String("/"),
		}

		allPlatforms := s3.NewListObjectsV2Paginator(client, platformParams)

		for allPlatforms.HasMorePages() {
			platsPage, platsErr := allPlatforms.NextPage(context.TODO())

			if platsErr != nil {
				log.Fatal(err)
				return []string{}
			}
			for _, platform := range platsPage.CommonPrefixes {
				fmt.Printf("..scanning %s\n", *platform.Prefix)

				platformParams := &s3.ListObjectsV2Input{
					Bucket:    aws.String(bucket),
					Prefix:    aws.String(*platform.Prefix),
					Delimiter: aws.String("/"),
				}

				platformPaginator := s3.NewListObjectsV2Paginator(client, platformParams)

				for platformPaginator.HasMorePages() {
					surveysPage, err := platformPaginator.NextPage(context.TODO())
					if err != nil {
						log.Fatal(err)
						return []string{}
					}

					for _, survey := range surveysPage.CommonPrefixes {
						fmt.Printf(" - %s\n", *survey.Prefix)
					}
				}
			}
		}
	}

	//for _, object := range result.Contents {
	//	log.Printf("key=%s size=%d", aws.ToString(object.Key), *object.Size)
	//}

	return surveyPaths
}
