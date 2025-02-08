package util

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	ErrCodeFailedToParseRequestBody = "failed-to-parse-request-body"
)

var (
	ErrFailedToParseRequestBody = errors.New("failed to parse request body: {error}")
)

func ConvertRequestBodyToMap(ctx context.Context, requestBody io.ReadCloser) (map[string]interface{}, *CustomError) {
	// Convert the request body to a map
	var body map[string]interface{}
	if err := json.NewDecoder(requestBody).Decode(&body); err != nil {
		cerr := NewCustomErrorWithKeys(ctx, ErrCodeFailedToParseRequestBody, ErrFailedToParseRequestBody, map[string]string{"error": err.Error()})
		return nil, cerr
	}
	return body, nil
}

func RespondWithError(ctx context.Context, w http.ResponseWriter, statusCode int, cerr *CustomError) {
	// Respond with an error message
	w.WriteHeader(statusCode)
	w.Write(cerr.GetErrorBytes(ctx))
}

func SendResponseMapWithStatus(ctx context.Context, w http.ResponseWriter, statusCode int, responseMap map[string]interface{}) {
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(responseMap)
	w.Write(response)
}
