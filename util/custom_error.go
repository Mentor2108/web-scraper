package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type CustomError struct {
	Message string
	Code    string
}

func (cerr *CustomError) Error() string {
	return fmt.Sprintf("| code: %s | error: %s", cerr.Code, cerr.Message)
}

func (cerr *CustomError) GetErrorMap(ctx context.Context) []byte {
	errorBytes, _ := json.Marshal(map[string]interface{}{
		"code":  cerr.Code,
		"error": cerr.Message,
	})
	return errorBytes
}

func NewCustomError(ctx context.Context, code string, err error) *CustomError {
	return &CustomError{
		Message: err.Error(),
		Code:    code,
	}
}

func NewCustomErrorWithKeys(ctx context.Context, code string, err error, keys map[string]string) *CustomError {
	errMessage := err.Error()
	for key, value := range keys {
		errMessage = strings.ReplaceAll(errMessage, "{"+key+"}", value)
	}

	return &CustomError{
		Message: errMessage,
		Code:    code,
	}
}
