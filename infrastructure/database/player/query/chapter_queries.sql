-- name: GetChapterProgressByPlayerID :many
SELECT
    chapter_id,
    checked_checkpoints
FROM
    chapter
WHERE
    player_id = $1;

-- name: GetChapterProgressByPlayerIDAndChapterID :one
SELECT
    chapter_id,
    checked_checkpoints
FROM
    chapter
WHERE
    player_id = $1 AND
    chapter_id = $2;

-- name: UpsertChapterProgressOnPlayer :one
INSERT INTO
    chapter (player_id, chapter_id, checked_checkpoints)
VALUES
    ($1, $2, $3)
ON CONFLICT (chapter_id, player_id) DO UPDATE
SET
    checked_checkpoints = EXCLUDED.checked_checkpoints,
    updated_at = NOW()
RETURNING
    chapter_id,
    checked_checkpoints;
