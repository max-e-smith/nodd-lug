package mb

import (
	"errors"
	"fmt"
	"github.com/max-e-smith/cruise-lug/internal/common"
	"github.com/spf13/viper"
	"log"
	"time"
)

var multibeamNODDBucket = "noaa-dcdb-bathymetry-pds" // https://noaa-dcdb-bathymetry-pds.s3.amazonaws.com/index.html

func MultibeamDownload(request MultibeamRequest) {
	request.resolveSurveys()
	request.checkDiskAvailability()
	request.downloadSurveys()

	if request.Error != nil {
		log.Fatal(request.Error)
	}
}

func (request *MultibeamRequest) checkDiskAvailability() {
	if request.Error != nil || len(request.Prefixes) == 0 {
		return
	}

	checkDisk := viper.GetBool("check")
	if !checkDisk {
		fmt.Println("Skipping disk space check.")
		return
	}

	bytes, estimateErr := common.GetDiskUsageEstimate(multibeamNODDBucket, request.S3Client, request.Prefixes)
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
	dryRun := viper.GetBool("try")
	if dryRun {
		fmt.Println("Skipping download due to dry run flag.")
		return
	}

	start := time.Now()
	defer logDownloadTime(start)

	order := common.Order{
		Bucket:      multibeamNODDBucket,
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

func logDownloadTime(start time.Time) {
	fmt.Printf("Download completed in %g hours.\n", common.HoursSince(start))
}
