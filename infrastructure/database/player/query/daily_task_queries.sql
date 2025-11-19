-- name: GetDailyTasksByPlayerID :many
SELECT
    task_id,
    progress,
    claimed,
    points_earned
FROM daily_task
WHERE player_id = $1
    AND day_start_at = $2;

-- name: IncreaseProgressForDailyTaskBatch :many
INSERT INTO daily_task (player_id, task_id, day_start_at, progress)
SELECT player_id, task_id, day_start_at, progress
FROM jsonb_to_recordset($1::jsonb) AS t(
    player_id UUID,
    task_id INT,
    day_start_at TIMESTAMPTZ,
    progress INT
)
ON CONFLICT (player_id, task_id, day_start_at)
DO UPDATE SET
    progress = daily_task.progress + EXCLUDED.progress,
    updated_at = NOW()
WHERE daily_task.progress + EXCLUDED.progress >= 0
RETURNING task_id,
          progress,
          claimed,
          points_earned;

-- name: ClaimDailyTask :one
INSERT INTO daily_task (
    player_id, task_id, day_start_at,
    claimed, points_earned)
VALUES ($1, $2, $3, TRUE, $4)
ON CONFLICT (player_id, task_id, day_start_at)
DO UPDATE SET
    claimed = TRUE,
    points_earned = EXCLUDED.points_earned,
    updated_at = NOW()
WHERE daily_task.claimed = FALSE
RETURNING task_id,
          progress,
          claimed,
          points_earned;
