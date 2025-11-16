package dto

type Team struct {
	Members  []TeamMember `db:"members"`
	TeamName string       `db:"team_name"`
}

type TeamMember struct {
	IsActive bool   `db:"is_active"`
	UserId   string `db:"user_id"`
	Username string `db:"username"`
}
