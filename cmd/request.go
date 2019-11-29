/*
Copyright © 2019 Brody Smith <brodygs9630@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type requestProperties struct {
	url              string
	method           string
	numberOfRequests int
	headers          map[string]string
}

type requestResult struct {
	statusCode int
	body       string
}

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Send an HTTP request to the specified URL.",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		method, _ := cmd.Flags().GetString("method")
		url, _ := cmd.Flags().GetString("url")
		amount, _ := cmd.Flags().GetInt("amount")

		headers, _ := cmd.Flags().GetString("headers")
		headersMap := make(map[string]string)
		if headers != "" {
			headersMap = prepareHeaders(headers)
		}

		filename, _ := cmd.Flags().GetString("output")
		if filename == "" {
			pid := os.Getpid()
			filename = fmt.Sprintf("./results-%d.log", pid)
		}

		properties := requestProperties{
			url:              url,
			method:           strings.ToUpper(method),
			numberOfRequests: amount,
			headers:          headersMap,
		}

		c := make(chan requestResult, properties.numberOfRequests)
		for i := 0; i < properties.numberOfRequests; i++ {
			wg.Add(1)
			go processRequest(properties, c, &wg)
		}
		wg.Wait()
		close(c)
		fmt.Println("Done processing requests...")

		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		var byteTotal int
		for v := range c {
			result := fmt.Sprintf("new_request\nStatus_Code: %d\nBody: %s\n\n", v.statusCode, v.body)

			bytes, err := f.WriteString(result)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				os.Exit(1)
			}
			byteTotal += bytes
		}
		fmt.Printf("wrote %d bytes\n", byteTotal)

		err = f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("method", "m", "get", "HTTP method to use for the request")
	requestCmd.Flags().StringP("url", "u", "", "URL to send the request to")
	requestCmd.Flags().IntP("amount", "n", 1, "Amount of requests to send")
	requestCmd.Flags().StringP("output", "o", "", "Path to file for results")
	requestCmd.Flags().StringP("headers", "H", "", "Header list formated as {key}:{value}, separated by commas")
}

// Submit request and send http.Response to channel 'c'.
func processRequest(properties requestProperties, c chan requestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Build the request
	req, err := http.NewRequest(properties.method, properties.url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	// Set Headers (if any)
	for key, value := range properties.headers {
		req.Header.Set(key, value)
	}

	// Build the HTTP Client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		c <- requestResult{
			statusCode: -1,
			body:       err.Error(),
		}
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		c <- requestResult{
			statusCode: resp.StatusCode,
			body:       string(body),
		}
	}
}

func prepareHeaders(input string) map[string]string {
	headerMap := make(map[string]string)
	temp := strings.ReplaceAll(input, " ", "")
	pairs := strings.Split(temp, ",")

	for _, v := range pairs {
		innerSlice := strings.Split(v, ":")
		if len(innerSlice) != 2 {
			fmt.Fprintf(os.Stderr, "The header %s is improperly formatted!", v)
			os.Exit(1)
		}
		headerMap[innerSlice[0]] = innerSlice[1]
	}
	return headerMap
}
