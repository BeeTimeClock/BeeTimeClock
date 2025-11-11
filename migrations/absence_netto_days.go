package migrations

import (
	"log"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
)

const MIGRATION_ABSENCE_NETTO_DAYS = "ABSENCE_NETTO_DAYS"

func MigrateAbsenceNettoDays(migrationRepo *repository.Migration, absenceRepo *repository.Absence, holidayRepo *repository.Holiday) error {
	_, err := migrationRepo.MigrationFindByTitle(MIGRATION_ABSENCE_NETTO_DAYS)
	migrationExists := true

	if err != nil {
		if err == repository.ErrMigrationNotFound {
			migrationExists = false
		} else {
			return err
		}
	}

	if migrationExists {
		log.Printf("Migration: %s already finished", MIGRATION_ABSENCE_NETTO_DAYS)
		return nil
	}

	log.Printf("Migration: %s started", MIGRATION_ABSENCE_NETTO_DAYS)
	absences, err := absenceRepo.FindAll(true)
	if err != nil {
		panic(err)
	}

	for _, absence := range absences {
		if absence.NettoDays != nil && *absence.NettoDays > 0 {
			continue
		}

		holidays, err := holidayRepo.HolidayFindByDateRange(absence.AbsenceFrom, absence.AbsenceTill)
		if err != nil {
			migration := model.Migration{
				Title:      MIGRATION_ABSENCE_NETTO_DAYS,
				Result:     err.Error(),
				FinishedAt: time.Now(),
				Success:    false,
			}
			migrationRepo.MigrationInsert(&migration)

			return err
		}

		absence.CalculateNettoDays(holidays)

		err = absenceRepo.Update(&absence)
		if err != nil {
			migration := model.Migration{
				Title:      MIGRATION_ABSENCE_NETTO_DAYS,
				Result:     err.Error(),
				FinishedAt: time.Now(),
				Success:    false,
			}
			migrationRepo.MigrationInsert(&migration)

			return err
		}
	}

	migration := model.Migration{
		Title:      MIGRATION_ABSENCE_NETTO_DAYS,
		Result:     "absence netto days migrated",
		FinishedAt: time.Now(),
		Success:    true,
	}
	migrationRepo.MigrationInsert(&migration)

	log.Printf("Migration: %s finished", MIGRATION_ABSENCE_NETTO_DAYS)
	return nil
}
