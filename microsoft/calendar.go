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

func CreateCalendarEntryFromAbsence(username string, absence *model.Absence) (string, error) {
	from := fmt.Sprintf("%sT00:00:00", absence.AbsenceFrom.Format(time.DateOnly))
	till := fmt.Sprintf("%sT00:00:00", absence.AbsenceTill.Add(24*time.Hour).Format(time.DateOnly))
	identifier := absence.Identifier.String()

	return CreateCalendarEntry(username, identifier, absence.AbsenceReason.Description, from, till, graphmodels.PRIVATE_SENSITIVITY)
}

func CreateCalendarEntryFromExternalWork(username string, externalWork *model.ExternalWork) (string, error) {
	from := externalWork.From.Format("2006-01-02T15:04:05")
	till := externalWork.Till.Format("2006-01-02T15:04:05")
	identifier := externalWork.Identifier.String()

	return CreateCalendarEntry(username, identifier, externalWork.Description, from, till, graphmodels.PRIVATE_SENSITIVITY)
}

func CreateCalendarEntry(username string, identifier string, description string, fromIso string, tillIso string, sensitivity graphmodels.Sensitivity) (string, error) {
	subject := fmt.Sprintf("BTC: %s", description)
	requestBody := graphmodels.NewEvent()
	requestBody.SetSubject(&subject)

	start := graphmodels.NewDateTimeTimeZone()
	start.SetDateTime(&fromIso)
	timeZone := "Europe/Berlin"
	start.SetTimeZone(&timeZone)
	requestBody.SetStart(start)

	end := graphmodels.NewDateTimeTimeZone()
	end.SetDateTime(&tillIso)
	end.SetTimeZone(&timeZone)
	requestBody.SetEnd(end)

	requestBody.SetTransactionId(&identifier)

	isAllDay := true
	requestBody.SetIsAllDay(&isAllDay)

	showAs := graphmodels.OOF_FREEBUSYSTATUS
	requestBody.SetShowAs(&showAs)

	reminder := false
	requestBody.SetIsReminderOn(&reminder)

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

func DeleteCalendarEntry(username string, externalEventId string) error {
	graphClient, err := getClient()
	if err != nil {
		return err
	}

	err = graphClient.Users().ByUserId(username).Events().ByEventId(externalEventId).Delete(context.Background(), nil)
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
