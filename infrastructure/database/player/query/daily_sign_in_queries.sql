-- name: GetDailySignInByPlayerID :one
SELECT
    id,
    week_start_at,
    claimed_days
FROM
    daily_sign_in
WHERE
    player_id = $1
    AND week_start_at = $2;

-- name: GetDailySignInByID :one
SELECT
    id,
    week_start_at,
    claimed_days
FROM
    daily_sign_in
WHERE
    player_id = $1
    AND id = $2;

-- name: InitDailySignIn :one
INSERT INTO
    daily_sign_in (player_id, week_start_at)
VALUES
    ($1, $2)
RETURNING
    id,
    week_start_at,
    claimed_days;

-- name: MarkDailySignInDays :one
UPDATE daily_sign_in
SET
    claimed_days = $3
WHERE
    player_id = $1
    AND id = $2
RETURNING
    id,
    week_start_at,
    claimed_days;
