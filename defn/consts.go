package defn

import "time"

const (
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 15 * time.Second
	ReadHeaderTimeout = 5 * time.Second
)

const (
	ContentTypeJSON       = "application/json"
	ContentTypePlainText  = "text/plain; charset=UTF-8"
	HTTPHeaderContentType = "Content-Type"
)
