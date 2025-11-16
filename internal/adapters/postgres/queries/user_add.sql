INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id)
DO UPDATE SET username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
