package response

type Response struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	Data          any    `json:"data,omitempty"`
}
