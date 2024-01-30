package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func compressString(input []byte) ([]byte, error) {
	var compressedBuffer bytes.Buffer
	writer := gzip.NewWriter(&compressedBuffer)

	_, err := writer.Write(input)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return compressedBuffer.Bytes(), nil
}

func PublishFlameGraph(flameFile string) error {
	file, err := os.Open(flameFile)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	compressedData, err := compressString(content)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(compressedData)
	fmt.Print(encoded)
	return nil
}
