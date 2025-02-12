package microsoft

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func GetMicrosoftTenantId() string {
	return os.Getenv("MICROSOFT_TENANT_ID")
}

func GetMicrosoftClientId() string {
	return os.Getenv("MICROSOFT_CLIENT_ID")
}

func GetMicrosoftClientSecret() string {
	return os.Getenv("MICROSOFT_CLIENT_SECRET")
}

func IsMicrosoftConnected() bool {
	return GetMicrosoftTenantId() != "" && GetMicrosoftClientId() != "" && GetMicrosoftClientSecret() != ""
}

func getClient() (*graph.GraphServiceClient, error) {
	oboCredential, err := azidentity.NewClientSecretCredential(GetMicrosoftTenantId(),
		GetMicrosoftClientId(), GetMicrosoftClientSecret(), nil)
	if err != nil {
		return nil, err
	}

	return graph.NewGraphServiceClientWithCredentials(oboCredential, []string{"https://graph.microsoft.com/.default"})
}

func CreateCalendarEntry(username string, absence *model.Absence) (string, error) {
	log.Printf("Microsoft Create Event: for %s from %s to %s", username, absence.AbsenceFrom, absence.AbsenceTill)
	subject := "BTC: Abwesend"
	requestBody := graphmodels.NewEvent()
	requestBody.SetSubject(&subject)

	fmt.Printf("Date: %s\n", fmt.Sprintf("%sT00:00:00", absence.AbsenceFrom.Format(time.DateOnly)))

	start := graphmodels.NewDateTimeTimeZone()
	dateTime := fmt.Sprintf("%sT00:00:00", absence.AbsenceFrom.Format(time.DateOnly))
	start.SetDateTime(&dateTime)
	timeZone := "Europe/Berlin"
	start.SetTimeZone(&timeZone)
	requestBody.SetStart(start)

	end := graphmodels.NewDateTimeTimeZone()
	enddateTime := fmt.Sprintf("%sT00:00:00", absence.AbsenceTill.Add(24*time.Hour).Format(time.DateOnly))
	end.SetDateTime(&enddateTime)
	end.SetTimeZone(&timeZone)
	requestBody.SetEnd(end)

	transactionId := absence.Identifier.String()
	requestBody.SetTransactionId(&transactionId)

	isAllDay := true
	requestBody.SetIsAllDay(&isAllDay)

	showAs := graphmodels.OOF_FREEBUSYSTATUS
	requestBody.SetShowAs(&showAs)

	reminder := false
	requestBody.SetIsReminderOn(&reminder)

	sensitivity := graphmodels.PRIVATE_SENSITIVITY
	requestBody.SetSensitivity(&sensitivity)

	graphClient, err := getClient()
	if err != nil {
		return "", err
	}

	result, err := graphClient.Users().ByUserId(username).Events().Post(context.Background(), requestBody, nil)
	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			log.Printf("OData error: %v", odataErr.Error())
			return "", err
		} else {
			log.Printf("Error getting user: %s", err)
			return "", err
		}
	}

	return *result.GetId(), nil
}

func DeleteCalendarEntry(username string, absence *model.Absence) error {
	graphClient, err := getClient()
	if err != nil {
		return err
	}

	err = graphClient.Users().ByUserId(username).Events().ByEventId(absence.ExternalEventID).Delete(context.Background(), nil)
	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			if *odataErr.GetErrorEscaped().GetCode() == "ErrorItemNotFound" {
				return nil
			}
			return fmt.Errorf("error getting event: %v", odataErr.GetErrorEscaped().GetMessage())
		}
		return fmt.Errorf("error getting event: %v", err)
	}
	return nil
}
