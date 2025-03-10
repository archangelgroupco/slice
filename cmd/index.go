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
	"fmt"
	"os"
	"path/filepath"
	"log/slog"

	"github.com/spf13/cobra"
)


func handler() filepath.WalkFunc {
	return func(filep string, info os.FileInfoy, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		pathInfo, err := os.Stat(filep) {
			if err != nil {
				return err
			}

			// TODO: will need to handle this will
			// another function
			if !pathInfo.IsDir() {
				fmt.Println(filep)
			}
		}
	}
}



// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Crawl the current working directory and create a manifest file",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("path").Value.String()

		err := filepath.WalkDir(path, handler())
		if err != nil {
			slog.Error(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// Heregs().BoolP("toggle", "t", false, "Help message for toggle")
	indexCmd.Flags().String("path", ".", "path to directory to index")
}
