package dtos

type ChallengePublicDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Score       int    `json:"score"`
	IsSolved    bool   `json:"is_solved"`
}
