package errorcode

import "fmt"

type AppError struct {
	Code    int
	Message string
	Details any `json:"details,omitempty"`
}

func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// 🔹 Authentication Errors (1xx)
var (
	ErrNoSession           = &AppError{Code: 101, Message: "No session found"}
	ErrInvalidToken        = &AppError{Code: 102, Message: "Invalid authentication token"}
	ErrUnsupportedProvider = &AppError{Code: 103, Message: "Unsupported provider"}
	ErrGettingSessionStore = &AppError{Code: 104, Message: "Error getting session store"}
	ErrSavingSessionStore  = &AppError{Code: 105, Message: "Error saving session store"}
	ErrInvalidOauthState   = &AppError{Code: 106, Message: "Invalid oauth state"}
	ErrExchangeToken       = &AppError{Code: 107, Message: "Error exchanging token"}
	ErrRefreshToken        = &AppError{Code: 108, Message: "Error refreshing token"}
	ErrFetchingOauthUser   = &AppError{Code: 109, Message: "Error fetching oauth user"}
)

// 🔹 Database Errors (2xx)
var (
	ErrTransactionStart  = &AppError{Code: 201, Message: "Failed to start transaction"}
	ErrTransactionCommit = &AppError{Code: 202, Message: "Transaction commit failed"}
	ErrScanningRow       = &AppError{Code: 203, Message: "Error scanning row"}
	ErrIteratingRows     = &AppError{Code: 204, Message: "Error iterating rows"}
	ErrDBQuery           = &AppError{Code: 205, Message: "Error querying database"}
	ErrDBCreate          = &AppError{Code: 206, Message: "Error creating data"}
	ErrDBUpdate          = &AppError{Code: 207, Message: "Error updating data"}
	ErrDBDelete          = &AppError{Code: 208, Message: "Error deleting data"}
)

// Functionality Errors (3xx)
var (
	ErrPhraseCreation    = &AppError{Code: 301, Message: "Phrase creation failed"}
	ErrPhraseTagCreation = &AppError{Code: 302, Message: "Phrase tag creation failed"}
	ErrManualSyncLimit   = &AppError{Code: 303, Message: "Manual sync limit reached"}
	ErrGuestIDCreation   = &AppError{Code: 304, Message: "Error creating guest ID"}
)

// User Errors (4xx)
var (
	ErrBadRequest       = &AppError{Code: 400, Message: "Bad request"}
	ErrNoPaginationData = &AppError{Code: 401, Message: "No pagination data given"}
	ErrNoSearchingData  = &AppError{Code: 402, Message: "No searching data given"}
)

// Other Errors (5xx)
var (
	ErrURLUnescape = &AppError{Code: 500, Message: "Error unescaping URL"}
)
