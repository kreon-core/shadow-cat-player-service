-- name: UpsertBattleHistory :one
INSERT INTO battle_history (player_id, game_mode, tower_id, floor, map_id, time_survived)
VALUES ($1, $2, $3, $4, $5, 0)
ON CONFLICT (player_id)
    WHERE completed_at IS NULL
DO UPDATE SET
    updated_at = NOW()
RETURNING *;

-- name: GetBattleHistoryByID :one
SELECT *
FROM battle_history
WHERE id = $1;

-- name: CompleteBattleHistory :one
UPDATE battle_history
SET completed_at = NOW(),
    time_survived = $2,
    monster_kills = $3,
    total_damage_dealt = $4,
    updated_at = NOW()
WHERE id = $1
    AND completed_at IS NULL
RETURNING *;

-- name: ExitBattleHistory :exec
UPDATE battle_history
SET exited_at = NOW(),
    updated_at = NOW()
WHERE id = $1
    AND completed_at IS NULL;
