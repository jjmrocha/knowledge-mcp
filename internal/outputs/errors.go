package outputs

type AppError struct {
	Message         string         `json:"error"`
	ErrorCode       string         `json:"error_code"`
	Details         map[string]any `json:"details"`
	SuggestedAction string         `json:"suggested_action"`
	Recoverable     bool           `json:"recoverable"`
}

func (e *AppError) Error() string {
	return e.ErrorCode
}
