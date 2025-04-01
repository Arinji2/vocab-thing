package errorcode

import "fmt"

type AppError struct {
	Code     int
	Message  string
	Readable string
	Details  any `json:"details,omitempty"`
}

func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		Code:     e.Code,
		Message:  e.Message,
		Readable: e.Readable,
		Details:  details,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// 🔹 Authentication Errors (1xx)
var (
	ErrNoSession           = &AppError{Code: 101, Message: "No session found", Readable: "Not Logged In"}
	ErrInvalidToken        = &AppError{Code: 102, Message: "Invalid authentication token", Readable: "Authentication failed"}
	ErrUnsupportedProvider = &AppError{Code: 103, Message: "Unsupported provider", Readable: "Authentication failed"}
	ErrGettingSessionStore = &AppError{Code: 104, Message: "Error getting session store", Readable: "Authentication failed"}
	ErrSavingSessionStore  = &AppError{Code: 105, Message: "Error saving session store", Readable: "Authentication failed"}
	ErrInvalidOauthState   = &AppError{Code: 106, Message: "Invalid oauth state", Readable: "Authentication failed"}
	ErrExchangeToken       = &AppError{Code: 107, Message: "Error exchanging token", Readable: "Authentication failed"}
	ErrRefreshToken        = &AppError{Code: 108, Message: "Error refreshing token", Readable: "Authentication failed"}
	ErrFetchingOauthUser   = &AppError{Code: 109, Message: "Error fetching oauth user", Readable: "Authentication failed"}
)

// 🔹 Database Errors (2xx)
var (
	ErrTransactionStart  = &AppError{Code: 201, Message: "Failed to start transaction", Readable: "Database operation failed"}
	ErrTransactionCommit = &AppError{Code: 202, Message: "Transaction commit failed", Readable: "Database operation failed"}
	ErrScanningRow       = &AppError{Code: 203, Message: "Error scanning row", Readable: "Database operation failed"}
	ErrIteratingRows     = &AppError{Code: 204, Message: "Error iterating rows", Readable: "Database operation failed"}
	ErrDBQuery           = &AppError{Code: 205, Message: "Error querying database", Readable: "Database operation failed"}
	ErrDBCreate          = &AppError{Code: 206, Message: "Error creating data", Readable: "Database operation failed"}
	ErrDBUpdate          = &AppError{Code: 207, Message: "Error updating data", Readable: "Database operation failed"}
	ErrDBDelete          = &AppError{Code: 208, Message: "Error deleting data", Readable: "Database operation failed"}
)

// Functionality Errors (3xx)
var (
	ErrPhraseCreation    = &AppError{Code: 301, Message: "Phrase creation failed", Readable: "Operation failed"}
	ErrPhraseTagCreation = &AppError{Code: 302, Message: "Phrase tag creation failed", Readable: "Operation failed"}
	ErrManualSyncLimit   = &AppError{Code: 303, Message: "Manual sync limit reached", Readable: "Limit reached"}
	ErrGuestIDCreation   = &AppError{Code: 304, Message: "Error creating guest ID", Readable: "Operation failed"}
)

// User Errors (4xx)
var (
	ErrBadRequest       = &AppError{Code: 400, Message: "Bad request", Readable: "Invalid input"}
	ErrNoPaginationData = &AppError{Code: 401, Message: "No pagination data given", Readable: "Invalid input"}
	ErrNoSearchingData  = &AppError{Code: 402, Message: "No searching data given", Readable: "Invalid input"}
)

// Other Errors (5xx)
var (
	ErrURLUnescape = &AppError{Code: 500, Message: "Error unescaping URL", Readable: "Operation failed"}
)
