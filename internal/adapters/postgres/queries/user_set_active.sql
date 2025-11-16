UPDATE users
SET is_active = $1
WHERE user_id = $2
RETURNING user_id, username, is_active, team_name;
