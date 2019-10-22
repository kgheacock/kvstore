package router

type ErrorForwardResponse struct {
	Error string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}