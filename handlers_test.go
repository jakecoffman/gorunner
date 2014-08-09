package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jakecoffman/gorunner/service"
)

func TestResponses(t *testing.T) {
	w := httptest.NewRecorder()

	uri := "/jobs"
	param := make(url.Values)

	r, err := http.NewRequest("GET", uri+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	jobList := service.NewJobList()
	c := ctx{jobList: jobList}

	status, _ := listJobs(c, w, r)

	if status != http.StatusOK {
		t.Errorf("Expected %v, %v; got %v, %v", http.StatusOK, status)
	}
}
