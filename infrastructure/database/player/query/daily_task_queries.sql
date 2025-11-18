-- name: GetDailyTasksByPlayerID :many
SELECT
    task_id,
    progress,
    claimed,
    points_earned
FROM daily_task
WHERE player_id = $1
    AND day_start_at = $2;

-- name: UpsertProgressForDailyTask :one
INSERT INTO daily_task (
    player_id,
    task_id,
    day_start_at,
    progress
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (player_id, task_id, day_start_at) DO UPDATE SET
    progress = EXCLUDED.progress,
    updated_at = NOW()
RETURNING task_id,
          progress,
          claimed,
          points_earned;

-- name: ClaimDailyTask :one
UPDATE daily_task
SET
    claimed = TRUE,
    points_earned = $4,
    updated_at = NOW()
WHERE player_id = $1
    AND task_id = $2
    AND day_start_at = $3
RETURNING task_id,
          progress,
          claimed,
          points_earned;
