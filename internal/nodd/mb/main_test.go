package mb

import (
	"testing"
)

func TestMultibeamRequest_CheckDiskAvailability(t *testing.T) {
	// TODO
}

func TestMultibeamRequest_Download(t *testing.T) {
	// TODO
}

func TestMultibeamRequest_ResolveSurveys(t *testing.T) {
	// TODO
}

func TestIsSurveyMatch(t *testing.T) {
	const expected = true
	surveys := []string{"survey1", "survey2", "survey3", "survey4", "survey5"}
	found := isSurveyMatch(surveys, "survey5")
	if found != expected {
		t.Errorf("isSurveyMatch returned %v, wanted %v", found, expected)
	}
}
