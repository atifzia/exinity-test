package response

// APIResponse is a standard response structure for the APIs
type APIResponse struct {
	StatusCode int                    `json:"status_code" xml:"status_code"`
	Message    string                 `json:"message" xml:"message"`
	Data       map[string]interface{} `json:"data,omitempty" xml:"data,omitempty"`
}
