package swu

import (
	alerts "github.com/opsgenie/opsgenie-go-sdk/alerts"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
)

type OpsGenieService interface {
	AlertOpsTeam(message string, description string, source string, allow_duplicates bool)
}

type opsGenieServiceImpl struct {
	apiKey string
}

func NewOpsGenieService(apiKey string) OpsGenieService {
	return opsGenieServiceImpl{
		apiKey: apiKey,
	}
}

func (s opsGenieServiceImpl) getAlertCli() *ogcli.OpsGenieAlertClient {
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(s.apiKey)
	alertCli, cliErr := cli.Alert()

	if cliErr != nil {
		panic(cliErr)
	}
	return alertCli
}

func (s opsGenieServiceImpl) alert(message string, description string, source string, recipients []string) {
	alertCli := s.getAlertCli()

	req := alerts.CreateAlertRequest{
		Message:     message,
		Description: description,
		Source:      source,
		Recipients:  recipients,
	}
	_, alertErr := alertCli.Create(req)
	if alertErr != nil {
		panic(alertErr)
	}
}

func (s opsGenieServiceImpl) hasOpenAlert(message string) bool {
	alertCli := s.getAlertCli()
	listreq := alerts.ListAlertsRequest{}
	listresp, listErr := alertCli.List(listreq)
	if listErr != nil {
		panic(listErr)
	}

	for _, alert := range listresp.Alerts {
		if alert.Message == message && alert.Status == "open" {
			return true
		}
	}
	return false
}

func (s opsGenieServiceImpl) AlertOpsTeam(message string, description string, source string, allow_duplicates bool) {
	if allow_duplicates || !s.hasOpenAlert(message) {
		s.alert(message, description, source, []string{"ops_team_schedule"})
	}
}
