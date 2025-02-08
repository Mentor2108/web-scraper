package defn

import "errors"

const (
	ErrCodeFailedToParseRequestBody       = "failed-to-parse-request-body"
	ErrCodeMissingRequiredField           = "missing-required-field"
	ErrCodeScrapePhaseLibraryNotSupported = "scrape-phase-libary-not-supported"
	ErrCodeInvalidDatabaseAction          = "invalid-database-action"
	ErrCodeChromedpError                  = "chromedp-error"
	ErrCodeScrapeUrlEmpty                 = "scrape-url-empty"
	ErrCodeDatabaseCreateOperationFailed  = "database-create-operation-failed"
	ErrCodeDatabaseUpdateOperationFailed  = "database-update-operation-failed"
	ErrCodeDatabaseGetOperationFailed     = "database-get-operation-failed"
	ErrCodeDatabaseDeleteOperationFailed  = "database-delete-operation-failed"
)

var (
	ErrFailedToParseRequestBody       = errors.New("failed to parse request body: {error}")
	ErrMissingRequiredField           = errors.New("missing required field: {required_field}")
	ErrScrapePhaseLibraryNotSupported = errors.New("given scrape-phase libary '{library}' is not supported")
	ErrInvalidDatabaseAction          = errors.New("invalid database action")
	ErrChromepError                   = errors.New("an error occured in chromedp: {error}")
	ErrScrapeUrlEmpty                 = errors.New("provided url to scrape is empty")
	ErrDatabaseCreateOperationFailed  = errors.New("database create action failed: {error}")
	ErrDatabaseUpdateOperationFailed  = errors.New("database update action failed: {error}")
	ErrDatabaseGetOperationFailed     = errors.New("database get action failed: {error}")
	ErrDatabaseDeleteOperationFailed  = errors.New("database delete action failed: {error}")
)
