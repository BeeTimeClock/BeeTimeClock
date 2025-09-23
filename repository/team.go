package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type Team struct {
	env *core.Environment
}

func NewTeam(env *core.Environment) *Team {
	return &Team{
		env: env,
	}
}

func (r *Team) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Team{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.TeamMember{})
	if err != nil {
		return err
	}

	return nil
}

var ErrTeamNotFound = errors.New("Team not found")

func (r Team) TeamFindAll(withData bool) ([]model.Team, error) {
	var items []model.Team
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	if withData {
		db = db.Preload("Members.User").Preload(clause.Associations)
	}

	result := db.Find(&items)

	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r Team) TeamsFindByUserId(userId uint) ([]model.Team, error) {
	var items []model.Team
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Where("id IN (?)", db.Table("beetc_team_member").Select("team_id").Where("user_id = ?", userId)).Preload("Members.User").Preload(clause.Associations).Find(&items)

	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r Team) TeamFindById(id uint, withData bool) (model.Team, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Team{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	if withData {
		db = db.Preload(clause.Associations)
	}

	var item model.Team
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Team{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Team{}, ErrTeamNotFound
	}
	return item, result.Error
}

func (r Team) TeamInsert(item *model.Team) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Team) TeamUpdate(item *model.Team) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Team) TeamDelete(item *model.Team) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

var ErrTeamMemberNotFound = errors.New("TeamMember not found")

func (r Team) TeamMemberFindAll() ([]model.TeamMember, error) {
	var items []model.TeamMember
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r Team) TeamMemberFindById(id uint) (model.TeamMember, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.TeamMember{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.TeamMember
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.TeamMember{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.TeamMember{}, ErrTeamMemberNotFound
	}
	return item, result.Error
}

func (r Team) TeamMemberInsert(item *model.TeamMember) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Team) TeamMemberUpdate(item *model.TeamMember) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Team) TeamMemberDelete(item *model.TeamMember) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r Team) TeamMemberFindByTeamId(teamId uint, withData bool) ([]model.TeamMember, error) {
	var items []model.TeamMember
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	if withData {
		db = db.Preload(clause.Associations)
	}

	result := db.Find(&items, "team_id = ?", teamId)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}
