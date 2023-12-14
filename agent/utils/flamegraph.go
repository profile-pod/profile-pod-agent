package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

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

	encoded := base64.StdEncoding.EncodeToString(content)
	fmt.Print(encoded)
	return nil
	// fgData := api.FlameGraphData{EncodedFile: encoded}

	// return api.PublishEvent(api.FlameGraph, fgData)
}
