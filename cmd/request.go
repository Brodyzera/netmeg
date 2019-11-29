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

		properties := requestProperties{
			url:              url,
			method:           strings.ToUpper(method),
			numberOfRequests: amount,
		}

		fmt.Println("request called")
		c := make(chan http.Response, properties.numberOfRequests)
		for i := 0; i < properties.numberOfRequests; i++ {
			wg.Add(1)
			go processRequest(properties, c, &wg)
		}
		wg.Wait()
		close(c)
		fmt.Println("Done")

		f, err := os.OpenFile("results.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		for v := range c {
			body, _ := ioutil.ReadAll(v.Body)
			v.Body.Close()
			result := fmt.Sprintf("new_request\nStatus_Code: %d\nBody: %s\n\n", v.StatusCode, body)

			bytes, err := f.WriteString(result)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				os.Exit(1)
			}
			fmt.Printf("wrote %d bytes\n", bytes)

			err = f.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("method", "m", "get", "HTTP method to use for the request")
	requestCmd.Flags().StringP("url", "u", "", "URL to send the request to")
	requestCmd.Flags().IntP("amount", "n", 1, "Amount of requests to send")
}

// Submit request and send http.Response to channel 'c'.
func processRequest(properties requestProperties, c chan http.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	// Build the request
	req, err := http.NewRequest(properties.method, properties.url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	// Build the HTTP Client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	c <- *resp
}
