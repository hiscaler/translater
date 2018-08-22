package response

type Error struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Success bool `json:"success"`
}

// {"success": false, "error": {"message": "error message"}}
type FailResponse struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

// {"success": true, "data": {"rawContent": "message", "content": "message"}}
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    SuccessData `json:"data"`
}

type SuccessData struct {
	RawContent string `json:"rawContent"`
	Content    string `json:"content"`
}
