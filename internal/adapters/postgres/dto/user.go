package dto

type User struct {
	IsActive bool   `db:"is_active"`
	TeamName string `db:"team_name"`
	UserId   string `db:"user_id"`
	Username string `db:"username"`
}
