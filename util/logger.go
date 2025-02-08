package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logger *log.Logger = nil

func InitiateGlobalLogger(ctx context.Context) *log.Logger {
	outputLocation := ctx.Value("output-format")
	output := ""
	var ok bool
	if output, ok = outputLocation.(string); !ok {
		output = ""
	}

	var outputWriter *os.File
	if strings.EqualFold(output, "") {
		outputWriter = os.Stdout
	} else {
		err := os.Mkdir(filepath.Dir(output), os.ModePerm)
		if err != nil {
			fmt.Printf("failed to initiate output logger at given location: %s\n", err.Error())
			os.Exit(1)
		}
		outputWriter, err = os.Create(output)
		if err != nil {
			fmt.Printf("failed to initiate output logger at given location: %s\n", err.Error())
			os.Exit(1)
		}
	}
	logger = log.New(outputWriter, "", log.Ldate|log.Ltime|log.Llongfile)
	return logger
}

func GetGlobalLogger(ctx context.Context) *log.Logger {
	return logger
}
