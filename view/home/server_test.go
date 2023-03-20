package home

import (
	"testing"
)

func TestGetTrendingList(t *testing.T) {
	data, err := GetTrendingList()
	if err != nil {
		t.Errorf("GetTrendingList returned an error: %v", err)
	}else{
		t.Logf("GetTrendingList returned: %v", data)
	}
}
