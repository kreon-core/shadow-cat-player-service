package dto

type ChapterProgress struct {
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	ChapterID          int          `json:"chapter_id"`
	CheckedCheckpoints map[int]bool `json:"checked_checkpoints"`
}
