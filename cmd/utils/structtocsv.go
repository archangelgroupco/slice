package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slice/internal/models"
)

func WriteToCSV(fileName string, data []models.Entry) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	// Write the header to the file
	err = writer.Write([]string{"relative_path", "file_extension", "mime_type", "parser_version"})
	if err != nil {
		log.Println(err)
	}

	// Write each row to the file
	for _, v := range data {
		row := []string{v.RelativePath,
			v.FileExtension,
			v.MimeType,
			fmt.Sprintf("%d", v.ParserVersion)}
		err := writer.Write(row)
		if err != nil {
			log.Println(err)
		}
	}

	writer.Flush()
	fmt.Println("csv file created")

	return nil
}
