package domain

type Team struct {
	Members  []TeamMember
	TeamName string
}

type TeamMember struct {
	IsActive bool
	UserId   string
	Username string
}
