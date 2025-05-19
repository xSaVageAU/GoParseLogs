package fileops

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// ReadFileContent reads the content of a file, automatically handling gzip compression if needed
func ReadFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var reader io.Reader = file

	// If file ends with .gz, use gzip reader
	if strings.HasSuffix(strings.ToLower(filePath), ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return "", err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
