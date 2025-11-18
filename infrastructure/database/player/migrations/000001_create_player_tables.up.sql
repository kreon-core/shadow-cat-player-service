CREATE TABLE player (
    id UUID PRIMARY KEY,
    level INT NOT NULL DEFAULT 0,
    exp INT NOT NULL DEFAULT 0,
    coins INT NOT NULL DEFAULT 0,
    gems INT NOT NULL DEFAULT 0,

    current_energy INT NOT NULL,
    max_energy INT NOT NULL,
    next_energy_at TIMESTAMPTZ,

    best_map JSONB NOT NULL,

    current_skin INT NOT NULL DEFAULT 0,
    equipped_props JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE skin (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    config_skin_id INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(player_id, config_skin_id)
);

CREATE TABLE prop (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    config_prop_id INT NOT NULL,
    level INT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(player_id, config_prop_id, level)
);

CREATE TABLE tower(
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    tower_id INT NOT NULL,
    ticket INT NOT NULL DEFAULT 0,
    highest_floor INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (tower_id, player_id)
);

CREATE TABLE chapter(
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    chapter_id INT NOT NULL,
    checked_checkpoints JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (chapter_id, player_id)
);


CREATE TABLE daily_sign_in(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    week_start_at TIMESTAMPTZ NOT NULL,
    claimed_days JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(player_id, week_start_at)
);

CREATE TABLE daily_task (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    task_id INT NOT NULL,
    day_start_at TIMESTAMPTZ NOT NULL,
    progress INT NOT NULL DEFAULT 0,
    claimed BOOLEAN NOT NULL DEFAULT FALSE,
    points_earned INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(player_id, task_id, day_start_at)
);

CREATE TABLE daily_shop(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    item INT NOT NULL,
    day_start_at TIMESTAMPTZ NOT NULL,
    quantity INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(player_id, item, day_start_at)
);

CREATE TABLE battle_history(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES player(id) ON UPDATE CASCADE ON DELETE CASCADE,

    game_mode INT NOT NULL,
    tower_id INT,
    floor INT,
    map_id INT,
    completed_at TIMESTAMPTZ,
    time_survived INT NOT NULL,
    monster_kills INT NOT NULL,
    total_damage_dealt INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
