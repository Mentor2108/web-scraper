package util

import (
	"backend-service/defn"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func ConvertRequestBodyToMap(ctx context.Context, requestBody io.ReadCloser) (map[string]interface{}, *CustomError) {
	// Convert the request body to a map
	var body map[string]interface{}
	if err := json.NewDecoder(requestBody).Decode(&body); err != nil {
		cerr := NewCustomErrorWithKeys(ctx, defn.ErrFailedToParseRequestBodyCode, defn.ErrFailedToParseRequestBody, map[string]string{"error": err.Error()})
		return nil, cerr
	}
	return body, nil
}

func RespondWithError(ctx context.Context, w http.ResponseWriter, statusCode int, cerr *CustomError) {
	// Respond with an error message
	w.WriteHeader(statusCode)
	w.Write(cerr.GetErrorMap(ctx))
}
