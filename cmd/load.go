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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

type requestJSON []struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Amount      int    `json:"amount"`
	Body        string `json:"body"`
	Headers     string `json:"headers"`
	Bfile       string `json:"bfile"`
	Hfile       string `json:"hfile"`
	Output      string `json:"output"`
	Mode        string `json:"mode"`
}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Execute request jobs stored in a JSON file.",
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("file")

		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		var requests requestJSON
		err = json.Unmarshal(bytes, &requests)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		for _, req := range requests {
			reqArgs := []string{"request", "-u", req.URL, "-m", req.Method, "-n", strconv.Itoa(req.Amount),
				"-b", req.Body, "-H", req.Headers, "--bfile", req.Bfile, "--hfile", req.Hfile, "-o", req.Output, "--mode", req.Mode}
			fmt.Println(reqArgs)

			out, err := exec.Command("netmeg", reqArgs...).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running request '%s': %s\n%s\n", req.Description, err, out)
				continue
			}
			fmt.Printf("%s\n", out)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().StringP("file", "f", "", "JSON file containing requests")
}
