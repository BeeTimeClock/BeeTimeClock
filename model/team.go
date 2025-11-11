package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Teamname string `gorm:"unique"`
	Members  []TeamMember
}

type TeamResponse struct {
	gorm.Model
	Teamname string
	Members  []TeamMemberResponse
}

func (t *Team) GetTeamResponse() TeamResponse {
	res := TeamResponse{
		Model:    t.Model,
		Teamname: t.Teamname,
		Members:  []TeamMemberResponse{},
	}

	for _, member := range t.Members {
		res.Members = append(res.Members, member.GetTeamMemberResponse())
	}

	return res
}

type TeamMemberResponse struct {
	gorm.Model
	UserID uint
	User   UserResponse
	Level  TeamLevel
}

const (
	TeamLevel_Lead          TeamLevel = "lead"
	TeamLevel_LeadSurrogate TeamLevel = "lead_surrogate"
	TeamLevel_Member        TeamLevel = "member"
)

type TeamLevel string

type TeamMember struct {
	gorm.Model
	TeamID uint `gorm:"index:idx_team_member,unique"`
	Team   Team
	UserID uint `gorm:"index:idx_team_member,unique"`
	User   User
	Level  TeamLevel
}

func (tm *TeamMember) GetTeamMemberResponse() TeamMemberResponse {
	return TeamMemberResponse{
		Model:  tm.Model,
		UserID: tm.UserID,
		User:   tm.User.GetUserResponse(),
		Level:  tm.Level,
	}
}

type TeamCreateRequest struct {
	Teamname   string `binding:"required"`
	TeamLeadId uint   `binding:"required"`
}

type TeamMemberCreateRequest struct {
	UserID uint      `binding:"required"`
	Level  TeamLevel `binding:"required"`
}
