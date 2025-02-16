package worker

import (
	"strings"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func NotifyAbsenceWeek(env *core.Environment, absenceRepo *repository.Absence) error {
	if !env.Notification.Enabled {
		return nil
	}

	now := time.Now()
	year, week := now.ISOWeek()

	weekStart := helper.WeekStart(year, week)
	weekEnd := weekStart.AddDate(0, 0, 5)

	absences, err := absenceRepo.FindByQuery(true, "absence_till >= ? and absence_till <= ?",
		weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	if err != nil {
		return err
	}

	mstClient := goteamsnotify.NewTeamsClient()
	grouped := make(map[time.Time][]string)
	var facts []adaptivecard.Fact

	for _, absence := range absences {
		for i := 0; i < 5; i++ {
			current := weekStart.AddDate(0, 0, i)
			if absence.IsDateInAbsence(current) {
				if _, exists := grouped[current]; !exists {
					grouped[current] = []string{}
				}

				grouped[current] = append(grouped[current], absence.User.FullName())
			}
		}
	}

	for weekday, names := range grouped {
		facts = append(facts, adaptivecard.Fact{
			Title: weekday.Weekday().String(),
			Value: strings.Join(names, ","),
		})
	}

	title := adaptivecard.NewTitleTextBlock("BTC: Abwesenheiten", false)
	factSet := adaptivecard.NewFactSet()
	factSet.AddFact(facts...)

	titleContainer := adaptivecard.NewContainer()
	titleContainer.AddElement(false, title)

	absenceContainer := adaptivecard.NewContainer()
	absenceContainer.AddElement(false, adaptivecard.Element(factSet))

	card := adaptivecard.NewCard()
	card.AddContainer(false, titleContainer)
	card.AddContainer(false, absenceContainer)

	msg := adaptivecard.NewMessage()
	msg.Attach(card)

	if err := mstClient.Send(env.Notification.WebhookUrl, msg); err != nil {
		return err
	}

	return nil
}
