-- name: GetInventoryByPlayerID :one
SELECT
    array_agg(skin.id) AS owned_skins,
    array_agg(
        json_build_object(
            'prop_id',
            prop.id,
            'config_prop_id',
            prop.config_prop_id,
            'level',
            prop.level,
            'quantity',
            prop.quantity
        )
    ) AS owned_props
FROM
    player
    INNER JOIN skin ON player.id = skin.player_id
    INNER JOIN prop ON player.id = prop.player_id
WHERE
    player.id = $1
GROUP BY
    player.id;

-- name: InsertOwnedSkins :many
INSERT INTO
    skin (player_id, config_skin_id)
SELECT
    player_id, config_skin_id
FROM jsonb_to_recordset($1::jsonb) AS p(
    player_id uuid,
    config_skin_id int
)
ON CONFLICT (player_id, config_skin_id) DO NOTHING
RETURNING skin.id;

-- name: RemoveOwnedSkins :many
DELETE FROM
    skin
WHERE
    player_id = $1
    AND config_skin_id = ANY($2::int[])
RETURNING id;

-- name: UpsertOwnedProps :many
INSERT INTO
    prop (player_id, config_prop_id, level, quantity)
SELECT
    player_id, config_prop_id, level, quantity
FROM jsonb_to_recordset($1::jsonb) AS p(
    player_id uuid,
    config_prop_id int,
    level int,
    quantity int
)
ON CONFLICT (player_id, config_prop_id, level) DO UPDATE
SET
    quantity = prop.quantity + EXCLUDED.quantity
RETURNING prop.id;

-- name: RemoveQuantityProp :one
UPDATE prop
SET
    quantity = quantity - $2,
    updated_at = NOW()
WHERE
    id = $1
    AND quantity >= $2
RETURNING id, quantity;
