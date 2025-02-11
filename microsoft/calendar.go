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
	subject := "BTC: Abwesend"
	requestBody := graphmodels.NewEvent()
	requestBody.SetSubject(&subject)

	start := graphmodels.NewDateTimeTimeZone()
	dateTime := absence.AbsenceFrom.Format(time.RFC3339)
	start.SetDateTime(&dateTime)
	timeZone := "Europe/Berlin"
	start.SetTimeZone(&timeZone)
	requestBody.SetStart(start)

	end := graphmodels.NewDateTimeTimeZone()
	enddateTime := absence.AbsenceTill.Format(time.RFC3339)
	end.SetDateTime(&enddateTime)
	endtimeZone := "Europe/Berlin"
	end.SetTimeZone(&endtimeZone)
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

	fmt.Printf("Delete Event: %s\n", absence.ExternalEventID)

	return graphClient.Users().ByUserId(username).Events().ByEventId(absence.ExternalEventID).Delete(context.Background(), nil)
}
