package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Teamname    string `gorm:"unique"`
	TeamOwnerID uint
	TeamOwner   User
	Members     []TeamMember
}

type TeamResponse struct {
	gorm.Model
	Teamname    string
	TeamOwnerID uint
	TeamOwner   UserResponse
	Members     []TeamMemberResponse
}

func (t *Team) GetTeamResponse() TeamResponse {
	res := TeamResponse{
		Model:       t.Model,
		TeamOwnerID: t.TeamOwnerID,
		Teamname:    t.Teamname,
		TeamOwner:   t.TeamOwner.GetUserResponse(),
		Members:     []TeamMemberResponse{},
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
}

type TeamMember struct {
	gorm.Model
	TeamID uint `gorm:"index:idx_team_member,unique"`
	Team   Team
	UserID uint `gorm:"index:idx_team_member,unique"`
	User   User
}

func (tm *TeamMember) GetTeamMemberResponse() TeamMemberResponse {
	return TeamMemberResponse{
		Model:  tm.Model,
		UserID: tm.UserID,
		User:   tm.User.GetUserResponse(),
	}
}

type TeamCreateRequest struct {
	Teamname    string `binding:"required"`
	TeamOwnerID uint   `binding:"required"`
}

type TeamMemberCreateRequest struct {
	UserID uint `binding:"required"`
}
