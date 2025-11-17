-- name: GetPlayerByID :one
SELECT
    id,
    level,
    exp,
    coins,
    gems,
    best_map,
    current_skin,
    equipped_props
FROM player
WHERE id = $1;

-- name: CreatePlayer :one
INSERT INTO player (
    id,
    coins,
    gems,
    current_energy,
    max_energy,
    best_map
) VALUES (
    $1, -- id UUID
    $2, -- coins INT
    $3, -- gems INT
    $4, -- current_energy INT
    $5, -- max_energy INT
    $6 -- best_map JSONB
)
RETURNING
    id,
    coins,
    gems,
    current_energy,
    max_energy,
    next_energy_at,
    best_map;
