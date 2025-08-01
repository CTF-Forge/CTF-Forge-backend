package dtos

// CreateChallengeRequestは問題作成APIのリクエストボディを定義します。
type CreateChallengeRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"` // カテゴリー名を文字列として受け取ります
	Score       int    `json:"score" binding:"required"`
	Flag        string `json:"flag" binding:"required"`
}

// ChallengeCreateResponseは問題作成成功時のレスポンスです。
type ChallengeCreateResponse struct {
	Message string `json:"message"`
}

// ErrorResponseはエラー発生時のレスポンスです。
type ErrorResponse struct {
	Error string `json:"error"`
}
