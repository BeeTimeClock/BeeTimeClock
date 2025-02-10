package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Teamname    string `gorm:"unique"`
	TeamOwnerID uint
	TeamOwner   User
}

type TeamMember struct {
	gorm.Model
	TeamID uint `gorm:"index:idx_team_member,unique"`
	Team   Team
	UserID uint `gorm:"index:idx_team_member,unique"`
	User   User
}

type TeamCreateRequest struct {
	Teamname    string `binding:"required"`
	TeamOwnerID uint   `binding:"required"`
}

type TeamMemberCreateRequest struct {
	UserID uint `binding:"required"`
}
