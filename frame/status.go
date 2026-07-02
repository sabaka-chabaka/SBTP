package frame

type Status uint16

const (
	StatusContinue  Status = 100
	StatusSwitching Status = 101
)

const (
	StatusOK        Status = 200
	StatusCreated   Status = 201
	StatusAccepted  Status = 202
	StatusNoContent Status = 204
)

const (
	StatusRedirect         Status = 300
	StatusMovedPermanently Status = 301
	StatusNotModified      Status = 304
)

const (
	StatusBadRequest       Status = 400
	StatusUnauthorized     Status = 401
	StatusForbidden        Status = 403
	StatusNotFound         Status = 404
	StatusMethodNotAllowed Status = 405
	StatusRequestTimeout   Status = 408
	StatusConflict         Status = 409
	StatusPayloadTooLarge  Status = 413
	StatusTooManyRequests  Status = 429
)

const (
	StatusInternalError      Status = 500
	StatusNotImplemented     Status = 501
	StatusBadGateway         Status = 502
	StatusServiceUnavailable Status = 503
	StatusGatewayTimeout     Status = 504
)

const (
	StatusInvalidMagic       Status = 600
	StatusUnsupportedVersion Status = 601
	StatusMalformedFrame     Status = 602
	StatusFrameTooLarge      Status = 603
	StatusChecksumMismatch   Status = 604
	StatusStreamClosed       Status = 605
	StatusCompressionError   Status = 606
)

var statusText = map[Status]string{
	StatusContinue:  "Continue",
	StatusSwitching: "Switching Protocols",

	StatusOK:        "OK",
	StatusCreated:   "Created",
	StatusAccepted:  "Accepted",
	StatusNoContent: "No Content",

	StatusRedirect:         "Redirect",
	StatusMovedPermanently: "Moved Permanently",
	StatusNotModified:      "Not Modified",

	StatusBadRequest:       "Bad Request",
	StatusUnauthorized:     "Unauthorized",
	StatusForbidden:        "Forbidden",
	StatusNotFound:         "Not Found",
	StatusMethodNotAllowed: "Method Not Allowed",
	StatusRequestTimeout:   "Request Timeout",
	StatusConflict:         "Conflict",
	StatusPayloadTooLarge:  "Payload Too Large",
	StatusTooManyRequests:  "Too Many Requests",

	StatusInternalError:      "Internal Server Error",
	StatusNotImplemented:     "Not Implemented",
	StatusBadGateway:         "Bad Gateway",
	StatusServiceUnavailable: "Service Unavailable",
	StatusGatewayTimeout:     "Gateway Timeout",

	StatusInvalidMagic:       "Invalid Magic",
	StatusUnsupportedVersion: "Unsupported Version",
	StatusMalformedFrame:     "Malformed Frame",
	StatusFrameTooLarge:      "Frame Too Large",
	StatusChecksumMismatch:   "Checksum Mismatch",
	StatusStreamClosed:       "Stream Closed",
	StatusCompressionError:   "Compression Error",
}

func (s Status) String() string {
	if text, ok := statusText[s]; ok {
		return text
	}
	return "Unknown Status"
}

func (s Status) IsSuccess() bool {
	return s >= 200 && s < 300
}

func (s Status) IsClientError() bool {
	return s >= 400 && s < 500
}

func (s Status) IsServerError() bool {
	return s >= 500 && s < 600
}

func (s Status) IsProtocolError() bool {
	return s >= 600 && s < 700
}

var errStatus = map[error]Status{
	ErrInvalidMagic:       StatusInvalidMagic,
	ErrUnsupportedVersion: StatusUnsupportedVersion,
	ErrMalformedFrame:     StatusMalformedFrame,
	ErrFrameTooLarge:      StatusFrameTooLarge,
	ErrMetadataTooLarge:   StatusFrameTooLarge,
	ErrChecksumMismatch:   StatusChecksumMismatch,
}

func StatusFromError(err error) Status {
	if status, ok := errStatus[err]; ok {
		return status
	}
	return StatusInternalError
}
