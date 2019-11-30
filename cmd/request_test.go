package cmd

import (
	"os"
	"testing"
)

func TestWriteResultsToFile(t *testing.T) {
	defer os.RemoveAll("./test_results")
	_ = os.Mkdir("./test_results", 0666)

	requestResults := []requestResult{
		requestResult{
			statusCode: 200,
			body:       "This is a test response body",
		},
		requestResult{
			statusCode: 200,
			body:       "This is another test response body",
		},
	}

	err := writeResultsToFile("./test_results/test-results-file.json", requestResults)
	if err != nil {
		t.Errorf("TestWriteResultsToFile failed, error occured: %s\n", err)
	}
}
