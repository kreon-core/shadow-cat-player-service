-- name: GetTowerProgressByPlayerID :many
SELECT
    tower_id,
    ticket,
    highest_floor
FROM
    tower
WHERE
    player_id = $1;

-- name: GetTowerProgressByPlayerIDAndTowerID :one
SELECT
    tower_id,
    ticket,
    highest_floor
FROM
    tower
WHERE
    player_id = $1
    AND tower_id = $2;

-- name: AddTowerTicketsToPlayer :one
INSERT INTO
    tower (player_id, tower_id, ticket)
VALUES
    ($1, $2, $3)
ON CONFLICT (tower_id, player_id) DO UPDATE
SET
    ticket = tower.ticket + EXCLUDED.ticket,
    updated_at = NOW()
RETURNING
    tower_id,
    ticket,
    highest_floor;

-- name: ConsumeTowerTicketsFromPlayer :one
UPDATE tower
SET
    ticket = ticket - $3,
    updated_at = NOW()
WHERE
    player_id = $1
    AND tower_id = $2
    AND ticket >= $3
RETURNING
    tower_id,
    ticket,
    highest_floor;

-- name: UpsertTowerProgressOnPlayer :exec
INSERT INTO
    tower (player_id, tower_id, highest_floor)
VALUES
    ($1, $2, $3)
ON CONFLICT (tower_id, player_id) DO UPDATE
SET
    highest_floor = GREATEST(tower.highest_floor, EXCLUDED.highest_floor),
    updated_at = NOW();
