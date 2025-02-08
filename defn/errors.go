package defn

import "errors"

const (
	ErrFailedToParseRequestBodyCode = "failed-to-parse-request-body"
	ErrMissingRequiredFieldCode     = "missing-required-field"
)

var (
	ErrFailedToParseRequestBody = errors.New("failed to parse request body: {error}")
	ErrMissingRequireField      = errors.New("missing required field: {required_field}")
)
