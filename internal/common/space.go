package common

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
	"log"
)

func getAvailableDiskSpace(localPath string) uint64 {
	usage := du.NewDiskUsage(localPath)
	if usage == nil {
		log.Fatalf("Could not get disk usage for path: %s", localPath)
	}
	return usage.Available() // bytes
}

func DiskSpaceCheck(neededBytes int64, targetPath string) error {
	fmt.Println("Checking available disk space")
	if neededBytes < 0 {
		neededBytes = 0
	}

	availableSpace := getAvailableDiskSpace(targetPath)

	fmt.Printf("  total disk size needed: %gGB\n", ByteToGB(neededBytes))
	fmt.Printf("  disk space available: %gGB\n", ByteToGB(int64(availableSpace)))

	if availableSpace > uint64(neededBytes) {
		return nil
	}

	return fmt.Errorf("not enough available space")
}
