-- name: GetInventoryByPlayerID :one
SELECT
    COALESCE(to_json(s.owned_skins), '[]') AS owned_skins,
    COALESCE(to_json(p.owned_props), '[]') AS owned_props
FROM player
LEFT JOIN (
    SELECT player_id, array_agg(config_skin_id) AS owned_skins
    FROM skin
    GROUP BY player_id
) s ON player.id = s.player_id
LEFT JOIN (
    SELECT player_id, array_agg(
        jsonb_build_object(
            'prop_id', id,
            'config_prop_id', config_prop_id,
            'level', level,
            'quantity', quantity
        )
    ) AS owned_props
    FROM prop
    GROUP BY player_id
) p ON player.id = p.player_id
WHERE player.id = $1;

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
