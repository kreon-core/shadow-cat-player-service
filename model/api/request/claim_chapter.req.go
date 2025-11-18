package request

type ClaimChapterRewards struct {
	ChapterID  int32 `json:"chapter_id"`
	Checkpoint int32 `json:"checkpoint"`
}
