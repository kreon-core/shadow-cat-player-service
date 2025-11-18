-- name: GetPlayerByID :one
SELECT
    id, level, exp, coins, gems,
    best_map,
    current_skin, equipped_props
FROM
    player
WHERE
    id = $1;

-- name: CreateNewPlayer :one
INSERT INTO
    player (
        id, coins, gems,
        current_energy, max_energy,
        best_map
    )
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    id, level, exp, coins, gems,
    best_map,
    current_skin, equipped_props;

-- name: UpdatePlayer :one
UPDATE player
SET
    level = $2,
    exp = $3,
    coins = $4,
    gems = $5,
    best_map = $6,
    current_skin = $7,
    equipped_props = $8,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    id, level, exp, coins, gems,
    best_map,
    current_skin, equipped_props;

-- name: GetPlayerEnergyByID :one
SELECT
    current_energy,
    max_energy,
    next_energy_at
FROM
    player
WHERE
    id = $1;

-- name: UpdatePlayerEnergy :one
UPDATE player
SET
    current_energy = $2,
    next_energy_at = $3,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    current_energy,
    max_energy,
    next_energy_at;
