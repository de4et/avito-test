SELECT u.user_id, u.username, u.is_active
FROM teams AS t
JOIN users AS u ON u.team_name = t.team_name
WHERE t.team_name = $1;
