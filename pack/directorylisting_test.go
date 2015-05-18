package lss

import (
	"os"
	"strings"
	"testing"
)

func Contains(list []string, value string) bool {
	for _, val := range list {
		if value == val {
			return true
		}
	}
	return false
}

func TestDirectorylisting_NoFilter(t *testing.T) {
	pth := os.Getenv("GOPATH") + "/tests/lsstest"
	err, ret := FilteredListingFromPath(pth, nil)
	if err != nil {
		t.Error(err)
	}
	validlist := []string{
		"foo.0001.mb",
		"foo.0002.mb",
		"foo.0003.mb",
		".foo.0004.mb",
	}
	cnt := 0
	for i, val := range ret {
		if Contains(validlist, val) != true {
			t.Error("FilteredListingFromPath did not contain", val)
		}
		cnt = i
	}
	if len(validlist) != cnt+1 {
		t.Error("Wrong number of items returned:", cnt, ".Should be:", len(validlist))
	}
}

func TestDirectorylisting_Filter(t *testing.T) {
	pth := os.Getenv("GOPATH") + "/tests/lsstest"
	err, ret := FilteredListingFromPath(pth, func(nm string) bool {
		if strings.Index(nm, ".") == 0 {
			return false
		}
		return true
	})
	if err != nil {
		t.Error(err)
	}
	validlist := []string{
		"foo.0001.mb",
		"foo.0002.mb",
		"foo.0003.mb",
	}
	cnt := 0
	for i, val := range ret {
		if Contains(validlist, val) != true {
			t.Error("FilteredListingFromPath did not contain", val)
		}
		cnt = i
	}
	if len(validlist) != cnt+1 {
		t.Error("Wrong number of items returned:", cnt, ".Should be:", len(validlist))
	}
}
