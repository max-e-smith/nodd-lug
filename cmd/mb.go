package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/max-e-smith/cruise-lug/internal/common"
	"github.com/max-e-smith/cruise-lug/internal/nodd/mb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var s3client s3.Client

var mbCmd = &cobra.Command{
	Use:   "mb",
	Short: "Handles multibeam bathymetry data requests",
	Long: `A cruise-lug command for downloading multibeam bathymetry
		   data. By default takes one or more survey names and a target
		   download directory as its arguments. By default will pulls
		   from the NODD (Noaa Open Data Dissemination) source. Usage:

			clug mb <options> <survey name> <survey name> <target directory>

			Options:
				-v --verbose (default: false)
					includes additional output in the console.
				-c --check (default: false)
					will check local disk space before downloading.
				-p --parallel <number> (default: 3)
					determines the number of parallel downloads for a request.
				-s --source <nodd | nccf> (default: nodd)
					determines the source of the multibeam data. Currently 
					only NODD is supported.
			`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath, surveys := parseArgs(cmd, args)
		parallelDownloads := getWorkersConfig()

		mb.MultibeamDownload(
			mb.MultibeamRequest{
				Surveys:     surveys,
				S3Client:    s3client,
				TargetDir:   targetPath,
				WorkerCount: parallelDownloads,
			},
		)

	},
}

func init() {
	RootCmd.AddCommand(mbCmd)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithRegion("us-east-1"),
	)

	if err != nil {
		fmt.Printf("Error loading AWS config: %s\n", err)
		fmt.Println("Failed to download multibeam surveys.")
		return
	}

	s3client = *s3.NewFromConfig(cfg)
}

func getWorkersConfig() int {
	numWorkers := viper.GetInt("parallel-downloads")
	if numWorkers < 1 {
		return 1
	}
	if numWorkers > 100 {
		return 100
	}
	return numWorkers
}

func parseArgs(cmd *cobra.Command, args []string) (string, []string) {
	var length = len(args)
	if length <= 1 {
		usageError(cmd, errors.New("please specify survey name(s) and a target file path"))
	}

	var targetPath = args[length-1]
	var surveys = args[:length-1]

	targetError := common.VerifyTargetPermissions(targetPath)
	if targetError != nil {
		usageError(cmd, targetError)
	}

	return targetPath, surveys
}

func usageError(cmd *cobra.Command, err error) {
	fmt.Println(cmd.UsageString())
	log.Fatal(err)
}
