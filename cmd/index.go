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
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"slice/internal/models"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var dirTreeIndex []models.Entry
var rootPath string

func handler() filepath.WalkFunc {
	return func(filep string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		ftype := mime.TypeByExtension(filepath.Ext(filep))

		pathInfo, err := os.Stat(filep)
		if err != nil {
			return err
		}

		if !pathInfo.IsDir() {
			relPath, err := filepath.Rel(rootPath, filep)
			if err != nil {
				log.Println(err)
			}

			entry := models.Entry{
				MimeType:      strings.Split(ftype, ";")[0],
				RelativePath:  relPath,
				FileExtension: filepath.Ext(filep),
				ParserVersion: 1,
			}
			dirTreeIndex = append(dirTreeIndex, entry)
		}

		return nil
	}
}

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Crawl the current working directory and create a manifest file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		path := cmd.Flag("path").Value.String()
		name := cmd.Flag("name").Value.String()

		// Set the root path so the handler func can determin the
		// root dir and only record relative
		rootPath = path

		err := filepath.Walk(path, handler())
		if err != nil {
			log.Println(err)
		}

		dsIndex := models.Manifest{
			DateTime: time.Now(),
			Name:     name,
			Nodes:    dirTreeIndex,
		}

		final, err := json.MarshalIndent(dsIndex, "", "	")
		if err != nil {
			log.Println(err)
		}
		fmt.Print(string(final))

	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// Heregs().BoolP("toggle", "t", false, "Help message for toggle")
	indexCmd.Flags().String("name", "manifest", "name of the manifest file")
	indexCmd.Flags().String("path", ".", "path to directory to index")
}
