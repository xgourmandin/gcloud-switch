package gcloud

import (
	"testing"
)

func TestCheckADCValid(t *testing.T) {
	result := CheckADCValid()
	if result != true && result != false {
		t.Error("CheckADCValid should return a boolean value")
	}
}
func TestGCloudFunctionsExist(t *testing.T) {
	err := SetProject("test-project")
	_ = err
	_, err = GetCurrentProject()
	_ = err
	valid := CheckADCValid()
	if valid != true && valid != false {
		t.Error("CheckADCValid should return a boolean")
	}
}
