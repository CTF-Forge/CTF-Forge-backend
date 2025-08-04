package dtos

// CreateChallengeRequestは問題作成APIのリクエストボディを定義します。
type CreateChallengeRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"` // カテゴリー名を文字列として受け取ります
	Score       int    `json:"score" binding:"required"`
	Flag        string `json:"flag" binding:"required"`
}

// UpdateChallengeRequest は問題更新APIのリクエストボディを定義します。
type UpdateChallengeRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Score       *int    `json:"score,omitempty"`
	Flag        *string `json:"flag,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

// ChallengeCreateResponseは問題作成成功時のレスポンスです。
type ChallengeCreateResponse struct {
	Message string `json:"message"`
}

// ChallengeDetailResponse は問題詳細取得APIのレスポンスです。
type ChallengeDetailResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    *string `json:"category"`
	Score       int     `json:"score"`
	Flag        string  `json:"flag"`
	IsPublic    bool    `json:"is_public"`
}

// ErrorResponseはエラー発生時のレスポンスです。
type ErrorResponse struct {
	Error string `json:"error"`
}
