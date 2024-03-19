package response

type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}
