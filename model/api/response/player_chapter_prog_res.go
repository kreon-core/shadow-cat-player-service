package response

type PlayerChapterProgress struct {
	Chapters []PlayerChapter `json:"chapters"`
}

type PlayerChapter struct {
	ChapterID          int          `json:"chapter_id"`
	CheckedCheckpoints map[int]bool `json:"checked_checkpoints"`
}
