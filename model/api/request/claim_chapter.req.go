package request

type ClaimChapterRewards struct {
	ChapterID  int32 `json:"chapter_id" validate:"required"`
	Checkpoint int32 `json:"checkpoint" validate:"required,min=0"`
}
