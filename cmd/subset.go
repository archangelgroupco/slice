/*
Copyright © 2025 archangelgroup.co

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"slice/cmd/utils"
	"slice/internal/models"

	"github.com/spf13/cobra"
)

// subsetCmd represents the subset command
var subsetCmd = &cobra.Command{
	Use:   "subset",
	Short: "create a data subset to be parsed by TheScribe",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		manifestFile := cmd.Flag("manifest-file").Value.String()
		outputFile := cmd.Flag("subset-file-name").Value.String()

		var manifest models.Manifest

		file, err := os.ReadFile(manifestFile)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(file, &manifest)
		if err != nil {
			log.Println(err)
		}

		err = utils.WriteToCSV(outputFile, manifest.Nodes)
	},
}

func init() {
	rootCmd.AddCommand(subsetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	subsetCmd.Flags().StringP("manifest-file", "f", "", "source manifest file")
	subsetCmd.Flags().StringP("subset-file-name", "o", "", "name of the output subset file")
	// subsetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
