package response

type Error struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Success bool `json:"success"`
}

// {'success': false, 'error': {'message': "error message"}}
type FailResponse struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    SuccessData `json:"data"`
}

type SuccessData struct {
	RawContent string `json:"rawContent"`
	Content    string `json:"content"`
}
