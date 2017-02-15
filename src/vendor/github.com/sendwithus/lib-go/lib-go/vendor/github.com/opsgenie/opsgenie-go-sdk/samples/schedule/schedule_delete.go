package main

import (
	"fmt"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
	"github.com/opsgenie/opsgenie-go-sdk/samples/constants"
	sch "github.com/opsgenie/opsgenie-go-sdk/schedule"
)

func main() {
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(constants.APIKey)

	schCli, cliErr := cli.Schedule()

	if cliErr != nil {
		panic(cliErr)
	}

	req := sch.DeleteScheduleRequest{Name: ""}
	response, schErr := schCli.Delete(req)

	if schErr != nil {
		panic(schErr)
	}

	fmt.Printf("status: %s\n", response.Status)
	fmt.Printf("code: %d\n", response.Code)
}
