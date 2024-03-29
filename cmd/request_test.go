package cmd

import (
	"os"
	"sync"
	"testing"
)

func TestWriteResultsToFile(t *testing.T) {
	defer os.RemoveAll("./test_results")
	_ = os.Mkdir("./test_results", 0777)

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

func TestWriteResultsToConsole(t *testing.T) {
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

	err := writeResultsToConsole(requestResults)
	if err != nil {
		t.Errorf("TestWriteResultsToConsole failed, error occured: %s\n", err)
	}
}

func TestPrepareHeaders(t *testing.T) {
	headersCorrect := "Content-Type:    application/json,     Test:another/header123"
	headersIncorrect := "Content-Type - application/json, Test:another/header123"
	_, errorCorrect := prepareHeaders(headersCorrect)
	if errorCorrect != nil {
		t.Errorf("TestPrepareHeaders failed, error occured: %s\n", errorCorrect)
	}

	_, errorIncorrect := prepareHeaders(headersIncorrect)
	if errorIncorrect == nil {
		t.Errorf("TestPrepareHeaders failed, expected header error didn't occur!\n")
	}
}

func TestProcessRequest(t *testing.T) {
	var wg sync.WaitGroup
	resultCount := 0
	headersMap := make(map[string]string)
	headersMap["Content-Type"] = "application/json"

	properties := requestProperties{
		url:              "https://localhost",
		method:           "get",
		numberOfRequests: 5,
		headers:          headersMap,
		body:             []byte(""),
	}
	c := make(chan requestResult, properties.numberOfRequests)

	for i := 0; i < properties.numberOfRequests; i++ {
		wg.Add(1)
		go processRequest(properties, c, &wg)
	}
	wg.Wait()
	close(c)

	for _ = range c {
		resultCount++
	}

	if resultCount != properties.numberOfRequests {
		t.Errorf("TestProcessRequest failed, received %d responses out of %d requests!", resultCount, properties.numberOfRequests)
	}
}
