package defn

import "errors"

const (
	ErrCodeFailedToParseRequestBody        = "failed-to-parse-request-body"
	ErrCodeMissingRequiredField            = "missing-required-field"
	ErrCodeScrapePhaseLibraryNotSupported  = "scrape-phase-libary-not-supported"
	ErrCodeProcessPhaseLibraryNotSupported = "process-phase-libary-not-supported"
	ErrCodeInvalidDatabaseAction           = "invalid-database-action"
	ErrCodeChromedpError                   = "chromedp-error"
	ErrCodeGoqueryError                    = "goquery-error"
	ErrCodeFailedToParseUrl                = "url-parse-failed"
	ErrCodeEmptyRootSelector               = "empty-root-selector"
	ErrCodeInvalidScrapingQueryContentType = "invalid-scraping-query-content-type"
	ErrCodeDatabaseCreateOperationFailed   = "database-create-operation-failed"
	ErrCodeDatabaseUpdateOperationFailed   = "database-update-operation-failed"
	ErrCodeDatabaseGetOperationFailed      = "database-get-operation-failed"
	ErrCodeDatabaseDeleteOperationFailed   = "database-delete-operation-failed"
	ErrCodeDirectoryCreationFailed         = "create-directory-failed"
	ErrCodeFileCreateFailed                = "file-create-failed"
	ErrCodeFileWriteFailed                 = "file-write-failed"
	ErrCodeFileStatFailed                  = "file-stat-failed"
	ErrCodeFileReadFailed                  = "file-read-failed"
	ErrCodeParseToIntFailed                = "parse-int-failed"
)

var (
	ErrFailedToParseRequestBody        = errors.New("failed to parse request body: {error}")
	ErrMissingRequiredField            = errors.New("missing required field: {field}")
	ErrScrapePhaseLibraryNotSupported  = errors.New("given scrape-phase libary '{library}' is not supported")
	ErrProcessPhaseLibraryNotSupported = errors.New("given process-phase libary '{library}' is not supported")
	ErrInvalidDatabaseAction           = errors.New("invalid database action")
	ErrChromedpError                   = errors.New("an error occured in chromedp: {error}")
	ErrGoqueryError                    = errors.New("an error occured in goquery: {error}")
	ErrFailedToParseUrl                = errors.New("an error occured while parsing url: {error}")
	ErrEmptyRootSelector               = errors.New("given root selector resulted in empty selection")
	ErrInvalidScrapingQueryContentType = errors.New("given scraping query content type '{type}' is invalid")
	ErrDatabaseCreateOperationFailed   = errors.New("database create action failed: {error}")
	ErrDatabaseUpdateOperationFailed   = errors.New("database update action failed: {error}")
	ErrDatabaseGetOperationFailed      = errors.New("database get action failed: {error}")
	ErrDatabaseDeleteOperationFailed   = errors.New("database delete action failed: {error}")
	ErrDirectoryCreationFailed         = errors.New("failed to create directory with path '{path}': {error}")
	ErrFileCreateFailed                = errors.New("failed to create file with path '{path}': {error}")
	ErrFileWriteFailed                 = errors.New("failed to write to file with path '{path}': {error}")
	ErrFileStatFailed                  = errors.New("failed to get file stat with path '{path}': {error}")
	ErrFileReadFailed                  = errors.New("failed to read file with path '{path}': {error}")
	ErrParseToIntFailed                = errors.New("failed to parse the field '{field}' to int: {error}")
)
